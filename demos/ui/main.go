package main

import (
	_ "embed"
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	imagex "github.com/moolmanruan/ebitengine-test/image"
	"github.com/moolmanruan/ebitengine-test/ui/panel"
	"image"
	_ "image/png"
	"log"
)

const (
	screenWidth  = 400
	screenHeight = 300
)

type Game struct {
	mainDialog *panel.T
	closeGame  bool // boolean indicating that the game should be closed
}

var ErrCloseGame = errors.New("close game")

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.mainDialog.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	p, err := loadPanel()
	if err != nil {
		log.Fatal(err)
	}

	panelGrid := imagex.NewGrid3x3(p, 49, 49, 49, 49)
	d := panel.New(panelGrid, image.Rect(50, 50, 350, 250))

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("UI Demo")
	if err := ebiten.RunGame(&Game{mainDialog: d}); err != nil {
		if !errors.Is(err, ErrCloseGame) {
			log.Fatal(err)
		}
	}
}
