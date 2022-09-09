package image

import (
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

// Grid represents a regular grid of images created from a bigger image.
type Grid struct {
	rows, cols int
	ii         []image.Image
}

func NewGrid(img *image.RGBA, x, y int) Grid {
	cols := img.Bounds().Size().X / x
	rows := img.Bounds().Size().Y / y

	ii := make([]image.Image, 0, cols*rows)
	for yi := 0; yi < rows; yi++ {
		for xi := 0; xi < cols; xi++ {
			xo, yo := xi*x, yi*y
			ii = append(ii, SubImage(img, image.Rect(xo, yo, xo+x, yo+y)))
		}
	}
	return Grid{rows: rows, cols: cols, ii: ii}
}

func NewGrid3x3(img image.Image, top, bottom, left, right int) Grid {
	eImg := ebiten.NewImageFromImage(img)
	w, h := eImg.Size()

	ii := make([]image.Image, 9)

	ii[0] = eImg.SubImage(image.Rect(0, 0, left, top))
	ii[1] = eImg.SubImage(image.Rect(left, 0, w-right, top))
	ii[2] = eImg.SubImage(image.Rect(w-right, 0, w, top))

	ii[3] = eImg.SubImage(image.Rect(0, top, left, h-bottom))
	ii[4] = eImg.SubImage(image.Rect(left, top, w-right, h-bottom))
	ii[5] = eImg.SubImage(image.Rect(w-right, top, w, h-bottom))

	ii[6] = eImg.SubImage(image.Rect(0, h-bottom, left, h))
	ii[7] = eImg.SubImage(image.Rect(left, h-bottom, w-right, h))
	ii[8] = eImg.SubImage(image.Rect(w-right, h-bottom, w, h))
	return Grid{rows: 3, cols: 3, ii: ii}
}

func (g Grid) List() []image.Image {
	return g.ii
}

func (g Grid) ImageAt(x, y int) image.Image {
	if x < 0 || x >= g.cols || y < 0 || y >= g.rows {
		return nil
	}
	return g.ii[x+y*g.cols]
}
