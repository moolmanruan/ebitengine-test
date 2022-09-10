package image

import (
	"image"
)

// SlicedImage represents an image sliced up into sub-images.
type SlicedImage struct {
	rows, cols int
	ii         []image.Image
}

// TileSlice splits a given image into sub-images with a width `w` and height
// `h` each.
func TileSlice(img image.Image, w, h int) SlicedImage {
	cols := img.Bounds().Size().X / w
	rows := img.Bounds().Size().Y / h

	ii := make([]image.Image, 0, cols*rows)
	for yi := 0; yi < rows; yi++ {
		for xi := 0; xi < cols; xi++ {
			xo, yo := xi*w, yi*h
			ii = append(ii, SubImage(img, image.Rect(xo, yo, xo+w, yo+h)))
		}
	}
	return SlicedImage{rows: rows, cols: cols, ii: ii}
}

// NineSlice splits a given image into 9 sub-image. This is used to draw the
// images at any size without distorting the corners and edged by keeping the
// corner tiles their normal size and stretching the inner tiles to fill the
// additional space.
func NineSlice(img image.Image, top, bottom, left, right int) SlicedImage {
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
	return SlicedImage{rows: 3, cols: 3, ii: ii}
}

// List returns all the sub-images images in a slice.
func (g SlicedImage) List() []image.Image {
	return g.ii
}

// RowSlice returns a set of `count` images in a specified row, with the first
// image being the one at `startCol`.
func (g SlicedImage) RowSlice(row, startCol, count int) []image.Image {
	ii := make([]image.Image, count)
	for ci := 0; ci < count; ci++ {
		ii[ci] = g.At(ci+startCol, row)
	}
	return ii
}

// ColSlice returns a set of `count` images in a specified column, with the
// first image being the one at `startRow`.
func (g SlicedImage) ColSlice(col, startRow, count int) []image.Image {
	ii := make([]image.Image, count)
	for ri := 0; ri < count; ri++ {
		ii[ri] = g.At(col, ri+startRow)
	}
	return ii
}

// At return the image at a specific coordinate
func (g SlicedImage) At(col, row int) image.Image {
	if col < 0 || col >= g.cols || row < 0 || row >= g.rows {
		return nil
	}
	return g.ii[col+row*g.cols]
}
