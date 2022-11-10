package byte2image

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
)

type Decoder interface {
	Decode(r io.Reader) (image.Image, error)
}

type PngDecoder struct{}

type JpegDecoder struct{}

type GeneralDecoder struct{}

func (d *PngDecoder) Decode(r io.Reader) (image.Image, error) {
	return png.Decode(r)
}

func (d *JpegDecoder) Decode(r io.Reader) (image.Image, error) {
	return jpeg.Decode(r)
}

func (d *GeneralDecoder) Decode(r io.Reader) (image.Image, error) {
	image, _, error := image.Decode(r)

	return image, error
}

func NewDecoder(url string) Decoder {
	if strings.Contains(url, ".png") {
		return &PngDecoder{}
	}
	if strings.Contains(url, ".jpg") {
		return &JpegDecoder{}
	}

	return &GeneralDecoder{}
}
