package main

import (
	_ "embed"
	imagex "github.com/moolmanruan/ebitengine-test/image"
	"image"
	_ "image/png"
)

//go:embed resources/panel.png
var imgBytes []byte

func loadPanel() (image.Image, error) {
	return imagex.FromBytes(imgBytes)
}
