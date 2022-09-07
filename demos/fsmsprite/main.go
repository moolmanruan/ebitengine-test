package main

import (
	_ "embed"
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/moolmanruan/ebitengine-test/sprite"
	"github.com/moolmanruan/ebitengine-test/sprite/fsm"
	_ "image/png"
	"log"
	"time"
)

const (
	screenWidth  = tileW * 2
	screenHeight = tileH * 2
)

type Game struct {
	closeGame bool // boolean indicating that the game should be closed
}

var ErrCloseGame = errors.New("close game")

func (g *Game) Update() error {
	if time.Now().Second()%2 == 0 {
		closeSprite.To(Normal)
	} else {
		closeSprite.To(Pressed)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	closeSprite.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

//go:embed resources/close.png
var imgBytes []byte

const tileW, tileH = 16, 16

var closeSprite *fsm.FSMSprite[ButtonState]

type ButtonState int

const (
	Normal  ButtonState = 0
	Pressed ButtonState = 1
)

func loadSprites() error {
	img, err := sprite.ImageFromBytes(imgBytes)
	if err != nil {
		return err
	}
	imageGrid := sprite.NewImageGrid(img, tileW, tileH)

	closeSprite = fsm.New[ButtonState]().
		AddState(Normal, imageGrid.ImageAt(0, 0)).
		AddState(Pressed, imageGrid.ImageAt(1, 0)).
		AddTransition(Normal, Pressed).
		AddTransition(Pressed, Normal).
		SetState(Normal)
	closeSprite.SetPosition(screenWidth/2-tileW/2, screenHeight/2-tileH/2)

	return nil
}

func main() {
	if err := loadSprites(); err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation Demo")
	if err := ebiten.RunGame(&Game{}); err != nil {
		if !errors.Is(err, ErrCloseGame) {
			log.Fatal(err)
		}
	}
}
