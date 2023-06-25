package main

import (
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

func ConvertToWebPImage(source image.Image, output io.Writer, quality float32) error {
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, quality)
	if err != nil {
		return err
	}

	return webp.Encode(output, source, options)
}

func ConvertToJPEGImage(source image.Image, output io.Writer, quality float32) error {
	return jpeg.Encode(output, source, &jpeg.Options{
		Quality: int(quality),
	})
}
