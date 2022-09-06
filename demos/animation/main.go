package main

import (
	_ "embed"
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/moolmanruan/ebitengine-test/sprite"
	"image"
	_ "image/png"
	"log"
	"time"
)

const (
	screenWidth  = runnerW * 3
	screenHeight = runnerH
)

type Game struct {
	closeGame bool // boolean indicating that the game should be closed
}

var ErrCloseGame = errors.New("close game")

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	runnerSprite.Draw(screen)
	idleSprite.Draw(screen)
	fallSprite.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

//go:embed resources/runner.png
var runnerBytes []byte

const runnerW, runnerH = 32, 32

var runnerSprite *sprite.AnimatedSprite
var idleSprite *sprite.AnimatedSprite
var fallSprite *sprite.AnimatedSprite

func loadSprites() error {
	img, err := sprite.ImageFromBytes(runnerBytes)
	if err != nil {
		return err
	}
	imageGrid := sprite.NewImageGrid(img, runnerW, runnerH)

	idleImages := []image.Image{
		imageGrid.ImageAt(0, 0),
		imageGrid.ImageAt(1, 0),
		imageGrid.ImageAt(2, 0),
		imageGrid.ImageAt(3, 0),
		imageGrid.ImageAt(4, 0),
	}
	runImages := []image.Image{
		imageGrid.ImageAt(0, 1),
		imageGrid.ImageAt(1, 1),
		imageGrid.ImageAt(2, 1),
		imageGrid.ImageAt(3, 1),
		imageGrid.ImageAt(4, 1),
		imageGrid.ImageAt(5, 1),
		imageGrid.ImageAt(6, 1),
		imageGrid.ImageAt(7, 1),
	}
	fallImages := []image.Image{
		imageGrid.ImageAt(0, 2),
		imageGrid.ImageAt(1, 2),
		imageGrid.ImageAt(2, 2),
		imageGrid.ImageAt(3, 2),
	}
	animationInterval := time.Millisecond * 100
	idleSprite = sprite.NewAnimated(idleImages, animationInterval)
	idleSprite.SetPosition(0, 0)
	runnerSprite = sprite.NewAnimated(runImages, animationInterval)
	runnerSprite.SetPosition(runnerW, 0)
	fallSprite = sprite.NewAnimated(fallImages, animationInterval)
	fallSprite.SetPosition(runnerW*2, 0)
	return nil
}

func main() {
	game := &Game{}

	if err := loadSprites(); err != nil {
		log.Fatal(err)
	}
	idleSprite.Play()
	runnerSprite.Play()
	fallSprite.Play()

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation Demo")
	if err := ebiten.RunGame(game); err != nil {
		if !errors.Is(err, ErrCloseGame) {
			log.Fatal(err)
		}
	}
}
