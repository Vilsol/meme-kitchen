package main

import (
	"bytes"
	"fmt"
	"github.com/golang/freetype/truetype"
	"github.com/llgcode/draw2d"
	"github.com/llgcode/draw2d/draw2dimg"
	"image"
	"image/color"
	"image/draw"
	"math"
	"memekitchen/data"
	"memekitchen/ent"
	"os"
	"path/filepath"
	"strings"
)

const RatioMeasureFontSize = 1000

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

func RenderMeme(payload []*data.Text, template *ent.Template, templateBytes []byte) (image.Image, error) {
	templateImage, _, err := image.Decode(bytes.NewReader(templateBytes))
	if err != nil {
		return nil, err
	}
	return RenderMemeImage(payload, template, templateImage)
}

func RenderMemeImage(payload []*data.Text, template *ent.Template, templateImage image.Image) (image.Image, error) {
	dest := image.NewRGBA(templateImage.Bounds())
	draw.Draw(dest, dest.Bounds(), templateImage, templateImage.Bounds().Min, draw.Src)

	// Initialize the graphic context on an RGBA image
	gc := draw2dimg.NewGraphicContext(dest)

	gc.SetFillColor(color.RGBA{R: 0x44, G: 0xff, B: 0x44, A: 0xff})
	gc.SetStrokeColor(color.RGBA{R: 0x44, G: 0x44, B: 0x44, A: 0xff})

	for i, t := range template.Data {
		if strings.TrimSpace(t.GetText()) == "" {
			continue
		}

		text := t.GetText()
		x := t.GetX()
		y := t.GetY()
		font := t.GetFont()
		fontSize := t.GetSize()
		fillColor := t.GetFillColor()
		strokeColor := t.GetStrokeColor()
		stroke := t.GetStroke()
		unfilled := t.GetUnfilled()
		horizontalAlign := t.GetHorizontalAlign()
		verticalAlign := t.GetVerticalAlign()
		width := t.GetWidth()
		height := t.GetHeight()

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

		gc.SetFontSize(float64(fontSize))
		gc.SetFontData(draw2d.FontData{Name: font, Family: draw2d.FontFamilySans, Style: draw2d.FontStyleNormal})

		//gc.SetStrokeColor(color.RGBA{
		//	R: 0xff,
		//	G: 0x00,
		//	B: 0x00,
		//	A: 0xff,
		//})
		//
		//gc.BeginPath()
		//gc.MoveTo(float64(x), float64(y))
		//gc.LineTo(float64(x+width), float64(y+height))
		//gc.Close()
		//gc.FillStroke()
		//
		//gc.BeginPath()
		//gc.MoveTo(float64(x+width), float64(y))
		//gc.LineTo(float64(x), float64(y+height))
		//gc.Close()
		//gc.FillStroke()
		//
		//gc.BeginPath()
		//gc.MoveTo(float64(x), float64(y))
		//gc.LineTo(float64(x+width), float64(y))
		//gc.LineTo(float64(x+width), float64(y+height))
		//gc.LineTo(float64(x), float64(y+height))
		//gc.LineTo(float64(x), float64(y))
		//gc.Close()
		//gc.Stroke()

		renderX := x
		switch horizontalAlign {
		case data.HorizontalAlign_LEFT:
			// Already left aligned
		case data.HorizontalAlign_CENTER:
			//renderX += width / 2
		case data.HorizontalAlign_RIGHT:
			//renderX += width
		}

		renderY := y
		switch verticalAlign {
		case data.VerticalAlign_TOP:
			// Already top aligned
		case data.VerticalAlign_MIDDLE:
			renderY += height / 2
		case data.VerticalAlign_BOTTOM:
			renderY += height
		}

		lines := strings.Split(strings.TrimSpace(text), "\n")

		widest := uint32(0)
		widestID := 0
		lineHeight := 1.35 * float64(fontSize)

		for j, l := range lines {
			left, _, right, _ := gc.GetStringBounds(l)
			if widest < uint32(right-left) {
				widest = uint32(right - left)
				widestID = j
			}
		}

		// Check width
		if widest > width {
			gc.SetFontSize(float64(RatioMeasureFontSize))
			left, _, right, _ := gc.GetStringBounds(lines[widestID])
			ratio := ((right - left) - float64(widest)) / (float64(RatioMeasureFontSize) - float64(fontSize))
			fontSize = uint32(math.Min(float64(fontSize), math.Floor(float64(width)/ratio)))
			gc.SetFontSize(float64(fontSize))
			lineHeight = 1.35 * float64(fontSize)
		}

		// Check height
		if float64(len(lines))*lineHeight > float64(height) {
			lineHeight = float64(height) / float64(len(lines))
			fontSize = uint32(math.Min(float64(fontSize), math.Floor(lineHeight/1.35)))
			gc.SetFontSize(float64(fontSize))
		}

		totalHeight := float64(len(lines)-1) * lineHeight

		switch verticalAlign {
		case data.VerticalAlign_TOP:
			// Do nothing
		case data.VerticalAlign_MIDDLE:
			renderY -= uint32(totalHeight / 2)
		case data.VerticalAlign_BOTTOM:
			renderY -= uint32(totalHeight)
		}

		if !unfilled && fillColor != "" {
			c, err := ParseHexColor(fillColor)
			if err != nil {
				return nil, err
			}

			gc.SetFillColor(c)
		}

		if stroke > 0 {
			// TODO Stroke size
			c, err := ParseHexColor(strokeColor)
			if err != nil {
				return nil, err
			}

			gc.SetLineWidth(float64(stroke))
			gc.SetStrokeColor(c)
		}

		for j, l := range lines {
			localX := renderX
			localY := renderY

			left, top, right, bottom := gc.GetStringBounds(l)

			switch horizontalAlign {
			case data.HorizontalAlign_CENTER:
				localX -= uint32((right - left) / 2)
				localX += width / 2
			case data.HorizontalAlign_LEFT:
				// Already left aligned
			case data.HorizontalAlign_RIGHT:
				localX -= uint32(right - left)
				localX += width
			}

			switch verticalAlign {
			case data.VerticalAlign_MIDDLE:
				localY += uint32((bottom - top) / 2)
			case data.VerticalAlign_TOP:
				localY += uint32(bottom - top)
			case data.VerticalAlign_BOTTOM:
			}

			if !t.Unfilled {
				gc.FillStringAt(l, float64(localX), float64(localY)+(float64(j)*lineHeight))
			}

			if t.Stroke != 0 {
				gc.StrokeStringAt(l, float64(localX), float64(localY)+(float64(j)*lineHeight))
			}
		}
	}

	return dest, nil
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
