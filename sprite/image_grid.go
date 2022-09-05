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

func (g ImageGrid) List() []image.Image {
	return g.ii
}

func (g ImageGrid) ImageAt(x, y int) (image.Image, error) {
	if x < 0 || x >= g.cols || y < 0 || y >= g.rows {
		return nil, nil
	}
	return g.ii[x+y*g.cols], nil
}
