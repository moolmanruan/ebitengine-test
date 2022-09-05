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

// ExtractSubImageTiles splits a give Image up and returns sub-images of the size
// of size (x,y).
func ExtractSubImageTiles(img image.Image, x, y int) []image.Image {
	eImg := ebiten.NewImageFromImage(img)
	xn := img.Bounds().Size().X / x
	yn := img.Bounds().Size().Y / y

	var ii []image.Image
	for xi := 0; xi < xn; xi++ {
		for yi := 0; yi < yn; yi++ {
			ii = append(ii, eImg.SubImage(image.Rect(xi*x, yi*y, (xi+1)*x, (yi+1)*y)))
		}
	}
	return ii
}
