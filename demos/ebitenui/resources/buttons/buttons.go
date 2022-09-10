package buttons

import (
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	imagex "github.com/moolmanruan/ebitengine-test/image"
	"image"
)

//go:embed button.png
var buttonBytes []byte

func NormalButton() *ebiten.Image {
	img, err := imagex.FromBytes(buttonBytes)
	if err != nil {
		return nil
	}
	return ebiten.NewImageFromImage(img.SubImage(image.Rect(0, 0, 49, 49)))
}

func PressedButton() *ebiten.Image {
	img, err := imagex.FromBytes(buttonBytes)
	if err != nil {
		return nil
	}
	return ebiten.NewImageFromImage(img.SubImage(image.Rect(49, 0, 98, 49)))
}

//go:embed panel.png
var panelBytes []byte

func Panel() *ebiten.Image {
	img, err := imagex.FromBytes(panelBytes)
	if err != nil {
		return nil
	}
	return ebiten.NewImageFromImage(img)
}
