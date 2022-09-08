package image

import (
	"bytes"
	"image"
	"image/draw"
)

func FromBytes(bb []byte) (*image.RGBA, error) {
	img, _, err := image.Decode(bytes.NewReader(bb))

	if imgRGBA, ok := img.(*image.RGBA); ok {
		return imgRGBA, nil
	}

	b := img.Bounds()
	newImg := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(newImg, newImg.Bounds(), img, b.Min, draw.Src)
	return newImg, err
}

func SubImage(src image.Image, rect image.Rectangle) *image.RGBA {
	dst := image.NewRGBA(image.Rect(0, 0, rect.Dx(), rect.Dy()))
	draw.Draw(dst, dst.Bounds(), src, rect.Min, draw.Src)
	return dst
}
