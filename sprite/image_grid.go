package sprite

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

func ImageFromBytes(bb []byte) (image.Image, error) {
	img, _, err := image.Decode(bytes.NewReader(bb))
	return img, err
}

// ImageGrid represents a regular grid of images created from a bigger image.
type ImageGrid struct {
	rows, cols int
	ii         []image.Image
}

func NewImageGrid(img image.Image, x, y int) ImageGrid {
	eImg := ebiten.NewImageFromImage(img)
	cols := img.Bounds().Size().X / x
	rows := img.Bounds().Size().Y / y

	ii := make([]image.Image, 0, cols*rows)
	for yi := 0; yi < rows; yi++ {
		for xi := 0; xi < cols; xi++ {
			xo, yo := xi*x, yi*y
			ii = append(ii, eImg.SubImage(image.Rect(xo, yo, xo+x, yo+y)))
		}
	}
	return ImageGrid{rows: rows, cols: cols, ii: ii}
}

func NewImageGrid3x3(img image.Image, top, bottom, left, right int) ImageGrid {
	eImg := ebiten.NewImageFromImage(img)
	w, h := eImg.Size()

	ii := make([]image.Image, 9)

	ii[0] = eImg.SubImage(image.Rect(0, 0, left, top))
	ii[1] = eImg.SubImage(image.Rect(left, 0, w-right, top))
	ii[2] = eImg.SubImage(image.Rect(w-right, 0, right, top))

	ii[3] = eImg.SubImage(image.Rect(0, top, left, h-bottom))
	ii[4] = eImg.SubImage(image.Rect(left, top, w-right, h-bottom))
	ii[5] = eImg.SubImage(image.Rect(w-right, top, right, h-bottom))

	ii[6] = eImg.SubImage(image.Rect(0, h-bottom, left, bottom))
	ii[7] = eImg.SubImage(image.Rect(left, h-bottom, w-right, bottom))
	ii[8] = eImg.SubImage(image.Rect(w-right, h-bottom, right, bottom))
	return ImageGrid{rows: 3, cols: 3, ii: ii}
}

func (g ImageGrid) List() []image.Image {
	return g.ii
}

func (g ImageGrid) ImageAt(x, y int) image.Image {
	if x < 0 || x >= g.cols || y < 0 || y >= g.rows {
		return nil
	}
	return g.ii[x+y*g.cols]
}
