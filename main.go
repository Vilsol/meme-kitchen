package main

import (
	"encoding/json"
	"fmt"
	"github.com/corona10/goimagehash"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	_ "github.com/kolesa-team/go-webp/webp"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"math/big"
	"memekitchen/config"
	"memekitchen/data"
	"memekitchen/ent"
	"memekitchen/ent/schema"
	template2 "memekitchen/ent/template"
	"memekitchen/storage"
)

type NewTemplate struct {
	Name string `form:"name"`
	Data string `form:"data"`
}

func main() {
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

	app := fiber.New(fiber.Config{
		StreamRequestBody: true,
	})

	app.Use(cors.New())

	app.Static("/", "./static")

	app.Get("/img/*.webp", func(c *fiber.Ctx) error {
		decoded, err := DecodeData(c.Params("*"))
		if err != nil {
			// TODO Debug
			log.Ctx(ctx).Err(err).Msg("failed decoding payload")
			return c.SendStatus(400)
		}

		// TODO Verify provided fonts exist
		// TODO Verify colors are hex

		template, err := db.Template.Get(c.Context(), int(decoded.Template))
		if err != nil {
			// TODO Debug
			return c.SendStatus(404)
		}

		templateFile, _, err := s3.Get(fmt.Sprintf("templates/images/%d.webp", template.ID))
		if err != nil {
			log.Ctx(ctx).Err(err).Msg("failed fetching template")

			// TODO Debug
			// TODO Error message
			return c.SendStatus(500)
		}

		meme, err := RenderMeme(decoded.Text, template, templateFile)
		if err != nil {
			log.Ctx(ctx).Err(err).Msg("failed rendering meme")

			// TODO Debug
			// TODO Error message
			return c.SendStatus(500)
		}

		// Cache for a week
		c.Response().Header.Set("Cache-Control", "public, max-age=604800")

		return c.Type("webp").Send(meme)
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

	app.Post("/api/templates", func(c *fiber.Ctx) error {
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

		// TODO Verify all params set for each text
		// TODO Verify provided fonts exist
		// TODO Verify colors are hex

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

		reader, writer := io.Pipe()

		var convertErr error
		go func() {
			defer writer.Close()
			convertErr = ConvertToWebPImage(img, writer, 100)
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

	if err := app.Listen(":3000"); err != nil {
		panic(err.Error())
	}
}
