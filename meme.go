package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"image/draw"
	"io"
	"memekitchen/data"
	"memekitchen/ent"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	dir, err := os.ReadDir("fonts")
	if err != nil {
		panic(err.Error())
	}

	for _, entry := range dir {
		if !entry.IsDir() {
			file, err := os.ReadFile(filepath.Join("fonts", entry.Name()))
			if err != nil {
				panic(err.Error())
			}

			fileName := entry.Name()[:strings.LastIndex(entry.Name(), ".")]

			ttf, err := truetype.Parse(file)
			if err != nil {
				panic(err.Error())
			}

			draw2d.RegisterFont(draw2d.FontData{Name: fileName, Family: draw2d.FontFamilySans, Style: draw2d.FontStyleNormal}, ttf)
		}
	}
}

func RenderMeme(payload []*data.Text, template *ent.Template, file io.ReadCloser) ([]byte, error) {
	templateImage, err := webp.Decode(file, nil)
	if err != nil {
		return nil, err
	}

	src, ok := templateImage.(*image.NRGBA)
	if !ok {
		return nil, errors.New("invalid template image")
	}

	dest := image.NewRGBA(src.Rect)
	draw.Draw(dest, dest.Bounds(), src, src.Bounds().Min, draw.Src)

	// Initialize the graphic context on an RGBA image
	gc := draw2dimg.NewGraphicContext(dest)

	gc.SetFillColor(color.RGBA{R: 0x44, G: 0xff, B: 0x44, A: 0xff})
	gc.SetStrokeColor(color.RGBA{R: 0x44, G: 0x44, B: 0x44, A: 0xff})

	for i, t := range template.Data {
		text := t.GetText()
		x := t.GetX()
		y := t.GetY()
		w := t.GetWidth()
		h := t.GetHeight()
		font := t.GetFont()
		fontSize := t.GetSize()
		fillColor := t.GetFillColor()
		strokeColor := t.GetStrokeColor()
		stroke := t.GetStroke()
		unfilled := t.GetUnfilled()
		horizontalAlign := t.GetHorizontalAlign()
		verticalAlign := t.GetVerticalAlign()

		t.GetWidth() // TODO Width

		if len(payload)-1 >= i {
			p := payload[i]

			if p.GetText() != "" {
				text = p.GetText()
			}

			if p.X != nil {
				x = p.GetX()
			}

			if p.Y != nil {
				y = p.GetY()
			}

			if p.Font != nil {
				font = p.GetFont()
			}

			if p.Size != nil {
				fontSize = p.GetSize()
			}

			if p.FillColor != nil {
				fillColor = p.GetFillColor()
			}

			if p.StrokeColor != nil {
				strokeColor = p.GetStrokeColor()
			}

			if p.Stroke != nil {
				stroke = p.GetStroke()
			}

			if p.Unfilled != nil {
				unfilled = p.GetUnfilled()
			}

			if p.HorizontalAlign != nil {
				horizontalAlign = p.GetHorizontalAlign()
			}

			if p.VerticalAlign != nil {
				verticalAlign = p.GetVerticalAlign()
			}
		}

		//gc.BeginPath()
		//gc.MoveTo(float64(x), float64(y))
		//gc.LineTo(float64(x+w), float64(y+h))
		//gc.Close()
		//gc.FillStroke()
		//
		//gc.BeginPath()
		//gc.MoveTo(float64(x+w), float64(y))
		//gc.LineTo(float64(x), float64(y+h))
		//gc.Close()
		//gc.FillStroke()

		gc.SetFontSize(float64(fontSize))
		gc.SetFontData(draw2d.FontData{Name: font, Family: draw2d.FontFamilySans, Style: draw2d.FontStyleNormal})

		left, top, right, bottom := gc.GetStringBounds(text)

		switch horizontalAlign {
		case data.HorizontalAlign_CENTER:
			x -= uint32((right - left) / 2)
			x += w / 2
		case data.HorizontalAlign_LEFT:
			// Already left aligned
		case data.HorizontalAlign_RIGHT:
			x -= uint32(right - left)
			x += w
		}

		switch verticalAlign {
		case data.VerticalAlign_MIDDLE:
			y += uint32((bottom - top) / 2)
			y += h / 2
		case data.VerticalAlign_TOP:
			y += uint32(bottom - top)
		case data.VerticalAlign_BOTTOM:
			y += h
		}

		if !unfilled {
			c, err := ParseHexColor(fillColor)
			if err != nil {
				return nil, err
			}

			gc.SetFillColor(c)
			gc.FillStringAt(text, float64(x), float64(y))
		}

		if stroke != 0 {
			// TODO Stroke size
			c, err := ParseHexColor(strokeColor)
			if err != nil {
				return nil, err
			}

			gc.SetStrokeColor(c)
			gc.StrokeStringAt(text, float64(x), float64(y))
		}
	}

	out := bytes.NewBuffer(make([]byte, 0))
	if err := ConvertToWebPImage(dest, out, 95); err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func ParseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4: %s", s)
	}
	return
}
