package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TwiN/go-away"
	"github.com/allegro/bigcache/v3"
	"github.com/corona10/goimagehash"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/etag"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/storage/redis/v2"
	_ "github.com/kolesa-team/go-webp/webp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	_ "image/png"
	"io"
	"math/big"
	"memekitchen/config"
	"memekitchen/data"
	"memekitchen/ent"
	"memekitchen/ent/schema"
	template2 "memekitchen/ent/template"
	"memekitchen/nsfw"
	"memekitchen/storage"
	"regexp"
	"strconv"
	"time"
)

type NewTemplate struct {
	Name string `form:"name"`
	Data string `form:"data"`
}

var extensionConversion = map[string]func(source image.Image, output io.Writer, quality float32, data string) error{
	"webp": ConvertToWebPImage,
	"jpeg": ConvertToJPEGImage,
	"jpg":  ConvertToJPEGImage,
}

var urlRegex = regexp.MustCompile(`[-a-zA-Z0-9@:%._+~#=]{1,256}\.[a-zA-Z0-9()]{1,6}\b([-a-zA-Z0-9()@:%_+.~#?&//=]*)`)

const maxTexts = 20

func ptr[T any](t T) *T {
	return &t
}

func main() {
	detector := nsfw.New("./nsfw_model")

	ctx := config.InitializeConfig()
	db := ConnectToDB()
	s3 := storage.ConnecToS3(storage.Config{
		Bucket:   viper.GetString("storage.bucket"),
		Key:      viper.GetString("storage.key"),
		Secret:   viper.GetString("storage.secret"),
		BaseURL:  viper.GetString("storage.base_url"),
		Endpoint: viper.GetString("storage.endpoint"),
		Region:   viper.GetString("storage.region"),
	})
	fiberStore := redis.New(redis.Config{
		Host:     viper.GetString("database.redis.host"),
		Port:     viper.GetInt("database.redis.port"),
		Password: viper.GetString("database.redis.pass"),
		Database: viper.GetInt("database.redis.db"),
	})

	cfg := bigcache.DefaultConfig(time.Hour)
	cfg.HardMaxCacheSize = 128
	cfg.MaxEntrySize = 2 << 19 // 1 MB

	templateCache, err := bigcache.New(ctx, cfg)
	if err != nil {
		panic(err)
	}

	app := fiber.New(fiber.Config{
		ProxyHeader:       fiber.HeaderXForwardedFor,
		StreamRequestBody: true,
	})

	app.Use(recover.New())

	app.Use(cors.New())

	app.Static("/", "./static", fiber.Static{
		Index: "index.html",
	})

	app.Use(NewLogger())

	renderer := app.Group("")

	renderer.Use(etag.New())

	renderer.Get("/img/:encoded.*", func(c *fiber.Ctx) error {
		extension := c.Params("*")
		if _, ok := extensionConversion[extension]; !ok {
			return c.SendStatus(400)
		}

		decoded, err := DecodeData(c.Params("encoded"))
		if err != nil {
			// TODO Debug
			log.Ctx(ctx).Err(err).Msg("failed decoding payload")
			return c.SendStatus(400)
		}

		// Verify max texts
		if len(decoded.Text) > maxTexts {
			// TODO Show an error image
			return c.SendStatus(400)
		}

		// Verify no URLs provided
		for _, text := range decoded.Text {
			if urlRegex.MatchString(text.GetText()) {
				// TODO Show an error image
				return c.SendStatus(400)
			}
		}

		// TODO Verify provided fonts exist
		// TODO Verify colors are hex

		template, err := db.Template.Get(c.Context(), int(decoded.Template))
		if err != nil {
			// TODO Debug
			return c.SendStatus(404)
		}

		templateBytes, err := templateCache.Get(strconv.Itoa(template.ID))
		if err != nil {
			if errors.Is(err, bigcache.ErrEntryNotFound) {
				templateFile, _, err := s3.Get(fmt.Sprintf("templates/images/%d.webp", template.ID))
				if err != nil {
					log.Ctx(ctx).Err(err).Msg("failed fetching template")

					// TODO Debug
					// TODO Error message
					return c.SendStatus(500)
				}

				templateBytes, err = io.ReadAll(templateFile)
				if err != nil {
					log.Ctx(ctx).Err(err).Msg("failed reading template")

					// TODO Debug
					// TODO Error message
					return c.SendStatus(500)
				}

				if err := templateCache.Set(strconv.Itoa(template.ID), templateBytes); err != nil {
					log.Ctx(ctx).Err(err).Msg("failed putting template in cache")
				}
			} else {
				log.Ctx(ctx).Err(err).Msg("failed reading cache")

				// TODO Debug
				// TODO Error message
				return c.SendStatus(500)
			}
		}

		meme, err := RenderMeme(decoded.Text, template, templateBytes)
		if err != nil {
			log.Ctx(ctx).Err(err).Msg("failed rendering meme")

			// TODO Debug
			// TODO Error message
			return c.SendStatus(500)
		}

		out := bytes.NewBuffer(make([]byte, 0))
		if err := extensionConversion[extension](meme, out, 95, c.Params("encoded")); err != nil {
			log.Ctx(ctx).Err(err).Msg("failed encoding meme to image")

			// TODO Debug
			// TODO Error message
			return c.SendStatus(500)
		}

		// Cache for a week
		c.Response().Header.Set("Cache-Control", "public, max-age=604800")

		return c.Type(extension).Send(out.Bytes())
	})

	app.Get("/api/templates", func(c *fiber.Ctx) error {
		// TODO Paging
		all, err := db.Template.Query().Select("id", "name").All(c.Context())
		if err != nil {
			// TODO Debug
			return c.SendStatus(500)
		}

		return c.JSON(all)
	})

	app.Get("/api/templates/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			// TODO Error message
			return c.SendStatus(400)
		}

		template, err := db.Template.Query().
			Select("id", "name", "data").
			Where(template2.ID(id)).
			First(c.Context())
		if err != nil {
			if ent.IsNotFound(err) {
				return c.SendStatus(404)
			}

			// TODO Error message
			return c.SendStatus(500)
		}

		return c.JSON(template)
	})

	app.Get("/template/:id.webp", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			// TODO Error message
			return c.SendStatus(400)
		}

		_, err = db.Template.Get(c.Context(), id)
		if err != nil {
			if ent.IsNotFound(err) {
				return c.SendStatus(404)
			}

			// TODO Error message
			return c.SendStatus(500)
		}

		templateFile, size, err := s3.Get(fmt.Sprintf("templates/images/%d.webp", id))
		if err != nil {
			// TODO Debug
			// TODO Error message
			return c.SendStatus(500)
		}

		// Cache for a week
		c.Response().Header.Set("Cache-Control", "public, max-age=604800")

		return c.Type("webp").SendStream(templateFile, int(size))
	})

	templateGroup := app.Group("")

	templateGroup.Use(limiter.New(limiter.Config{
		Max:        5,
		Expiration: time.Minute,
		Storage:    fiberStore,
	}))

	templateGroup.Post("/api/templates", func(c *fiber.Ctx) error {
		newTemplate := &NewTemplate{}
		if err := c.BodyParser(newTemplate); err != nil {
			// TODO Error message
			return c.SendStatus(400)
		}

		textData := make([]data.TemplateText, 0)
		if err := json.Unmarshal([]byte(newTemplate.Data), &textData); err != nil {
			// TODO Error message
			return c.SendStatus(400)
		}

		for _, text := range textData {
			// Verify no URLs provided
			if urlRegex.MatchString(text.GetText()) {
				// TODO Error message
				return c.SendStatus(400)
			}

			// Verify no profanities in template
			if goaway.IsProfane(text.GetText()) {
				// TODO Error message
				return c.SendStatus(400)
			}
		}

		imgMeta, err := c.FormFile("image")
		if err != nil {
			// TODO Error message
			return c.SendStatus(400)
		}

		imgFile, err := imgMeta.Open()
		if err != nil {
			// TODO Error message
			return c.SendStatus(400)
		}

		img, _, err := image.Decode(imgFile)
		if err != nil {
			// TODO Error message
			return c.SendStatus(400)
		}

		imageBuffer := bytes.NewBuffer(make([]byte, 0))
		if err := png.Encode(imageBuffer, img); err != nil {
			panic(err)
		}

		labels, err := detector.Labels(imageBuffer.Bytes())
		if err != nil {
			panic(err)
		}

		if !labels.IsSafe() {
			// TODO Error message
			return c.SendStatus(500)
		}

		avgDist, err := goimagehash.ExtAverageHash(img, 8, 8)
		if err != nil {
			// TODO Error message
			return c.SendStatus(500)
		}

		diffDist, err := goimagehash.ExtDifferenceHash(img, 8, 8)
		if err != nil {
			// TODO Error message
			return c.SendStatus(500)
		}

		perceptionDist, err := goimagehash.ExtPerceptionHash(img, 8, 8)
		if err != nil {
			// TODO Error message
			return c.SendStatus(500)
		}

		avgDistBig := schema.BigInt(*big.NewInt(0).SetUint64(avgDist.GetHash()[0]))
		diffDistBig := schema.BigInt(*big.NewInt(0).SetUint64(diffDist.GetHash()[0]))
		perceptionDistBig := schema.BigInt(*big.NewInt(0).SetUint64(perceptionDist.GetHash()[0]))

		// TODO check if DB already contains similar meme

		tx, err := db.Tx(c.Context())
		if err != nil {
			// TODO Error message
			return c.SendStatus(500)
		}

		template, err := tx.Template.Create().
			SetName(newTemplate.Name).
			SetAvgDistance(&avgDistBig).
			SetDiffDistance(&diffDistBig).
			SetPerceptionDistance(&perceptionDistBig).
			SetData(textData).
			Save(c.Context())
		if err != nil {
			// TODO Error message
			return c.SendStatus(500)
		}

		rewritten := make([]*data.Text, len(textData))
		for i, d := range textData {
			b := d
			rewritten[i] = &data.Text{
				TemplateText:    ptr(uint32(i)),
				Text:            b.Text,
				X:               &b.X,
				Y:               &b.Y,
				Width:           &b.Width,
				Height:          &b.Height,
				Font:            &b.Font,
				Size:            &b.Size,
				Unfilled:        &b.Unfilled,
				FillColor:       &b.FillColor,
				StrokeColor:     &b.StrokeColor,
				Stroke:          &b.Stroke,
				HorizontalAlign: &b.HorizontalAlign,
				VerticalAlign:   &b.VerticalAlign,
			}
		}

		_, err = RenderMemeImage(rewritten, template, img)
		if err != nil {
			_ = tx.Rollback()
			// TODO Error message
			return c.SendStatus(400)
		}

		reader, writer := io.Pipe()

		var convertErr error
		go func() {
			defer writer.Close()
			convertErr = ConvertToWebPImage(img, writer, 100, "")
		}()

		if _, err := s3.Put(c.Context(), fmt.Sprintf("templates/images/%d.webp", template.ID), reader); err != nil {
			_ = tx.Rollback()

			log.Ctx(ctx).Err(err).Msg("failed uploading template")

			// TODO Error message
			return c.SendStatus(500)
		}

		if convertErr != nil {
			// TODO Error message
			return c.SendStatus(500)
		}

		if err := tx.Commit(); err != nil {
			// TODO Error message
			return c.SendStatus(500)
		}

		return c.JSON(template)
	})

	app.Get("*", func(c *fiber.Ctx) error {
		return c.SendFile("static/200.html")
	})

	if err := app.Listen(":3000"); err != nil {
		panic(err.Error())
	}
}
