package image

import (
	"image"
)

// Grid represents a regular grid of images created from a bigger image.
type Grid struct {
	rows, cols int
	ii         []image.Image
}

// TileSlice splits a given image into sub-images with a width `w` and height
// `h` each.
func TileSlice(img image.Image, w, h int) Grid {
	cols := img.Bounds().Size().X / w
	rows := img.Bounds().Size().Y / h

	ii := make([]image.Image, 0, cols*rows)
	for yi := 0; yi < rows; yi++ {
		for xi := 0; xi < cols; xi++ {
			xo, yo := xi*w, yi*h
			ii = append(ii, SubImage(img, image.Rect(xo, yo, xo+w, yo+h)))
		}
	}
	return Grid{rows: rows, cols: cols, ii: ii}
}

// NineSlice splits a given image into 9 sub-image. This is used to draw the
// images at any size without distorting the corners and edged by keeping the
// corner tiles their normal size and stretching the inner tiles to fill the
// additional space.
func NineSlice(img image.Image, top, bottom, left, right int) Grid {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	ii := make([]image.Image, 9)

	ii[0] = SubImage(img, image.Rect(0, 0, left, top))
	ii[1] = SubImage(img, image.Rect(left, 0, w-right, top))
	ii[2] = SubImage(img, image.Rect(w-right, 0, w, top))

	ii[3] = SubImage(img, image.Rect(0, top, left, h-bottom))
	ii[4] = SubImage(img, image.Rect(left, top, w-right, h-bottom))
	ii[5] = SubImage(img, image.Rect(w-right, top, w, h-bottom))

	ii[6] = SubImage(img, image.Rect(0, h-bottom, left, h))
	ii[7] = SubImage(img, image.Rect(left, h-bottom, w-right, h))
	ii[8] = SubImage(img, image.Rect(w-right, h-bottom, w, h))
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
