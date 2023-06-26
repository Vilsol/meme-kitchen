package main

import (
	"github.com/dsoprea/go-exif/v3"
	exifcommon "github.com/dsoprea/go-exif/v3/common"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/jpeg"
	_ "image/png"
	"io"
)

func ConvertToWebPImage(source image.Image, output io.Writer, quality float32, data string) error {
	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, quality)
	if err != nil {
		return err
	}

	err = webp.Encode(output, source, options)
	if err != nil {
		return err
	}

	if len(data) > 0 {
		b, err := encodeExif(data)
		if err != nil {
			return err
		}

		if _, err := output.Write(b); err != nil {
			return err
		}
	}

	return nil
}

func ConvertToJPEGImage(source image.Image, output io.Writer, quality float32, data string) error {
	err := jpeg.Encode(output, source, &jpeg.Options{
		Quality: int(quality),
	})
	if err != nil {
		return err
	}

	if len(data) > 0 {
		b, err := encodeExif(data)
		if err != nil {
			return err
		}

		if _, err := output.Write(b); err != nil {
			return err
		}
	}

	return nil
}

func encodeExif(data string) ([]byte, error) {
	im, err := exifcommon.NewIfdMappingWithStandard()
	if err != nil {
		return nil, err
	}

	ti := exif.NewTagIndex()
	ib := exif.NewIfdBuilder(im, ti, exifcommon.IfdStandardIfdIdentity, exifcommon.TestDefaultByteOrder)

	ve := exifcommon.NewValueEncoder(exifcommon.EncodeDefaultByteOrder)
	encoded, err := ve.Encode(data)
	if err != nil {
		return nil, err
	}

	bt := exif.NewBuilderTag(ib.IfdIdentity().Name(), ^uint16(0), encoded.Type, exif.NewIfdBuilderTagValueFromBytes(encoded.Encoded), exifcommon.EncodeDefaultByteOrder)

	if err := ib.Add(bt); err != nil {
		return nil, err
	}

	return exif.NewIfdByteEncoder().EncodeToExif(ib)
}
