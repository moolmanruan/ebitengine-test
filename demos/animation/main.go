package main

import (
	_ "embed"
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/moolmanruan/ebitengine-test/animate"
	imagex "github.com/moolmanruan/ebitengine-test/image"
	"github.com/moolmanruan/ebitengine-test/sprite"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"time"
)

const (
	runnerW, runnerH = 32, 32
	progressBarH     = 10
	screenWidth      = runnerW * 3
	screenHeight     = runnerH + progressBarH
)

type Game struct {
	closeGame bool // boolean indicating that the game should be closed
}

var ErrCloseGame = errors.New("close game")

func (g *Game) Update() error {
	return nil
}

var (
	emptyImage    = ebiten.NewImage(3, 3)
	emptySubImage = emptyImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)
)

func init() {
	emptyImage.Fill(color.White)
}

func drawProgressBar(screen *ebiten.Image, x, y, progress float64, rgb [3]float32) {
	var path vector.Path

	xl, xr := float32(x), float32(x+runnerW*progress)
	yt, yb := float32(y), float32(y+runnerH)
	path.MoveTo(xl, yt)
	path.LineTo(xr, yt)
	path.LineTo(xr, yb)
	path.LineTo(xl, yb)
	path.LineTo(xl, yt)
	op := &ebiten.DrawTrianglesOptions{
		FillRule: ebiten.EvenOdd,
	}
	vs, is := path.AppendVerticesAndIndicesForFilling(nil, nil)
	for i := range vs {
		vs[i].SrcX = 1
		vs[i].SrcY = 1
		vs[i].ColorR = rgb[0]
		vs[i].ColorG = rgb[1]
		vs[i].ColorB = rgb[2]
	}

	screen.DrawTriangles(vs, is, emptySubImage, op)
}

func colorFromHex(v int) [3]float32 {
	return [3]float32{
		float32((v>>16)%0x100) / 0xff,
		float32((v>>8)%0x100) / 0xff,
		float32(v%0x100) / 0xff}
}

func (g *Game) Draw(screen *ebiten.Image) {
	runnerSprite.Draw(screen)
	idleSprite.Draw(screen)
	fallSprite.Draw(screen)
	drawProgressBar(screen, 0, runnerH, idleP, colorFromHex(0xfe4a49))
	drawProgressBar(screen, runnerW, runnerH, runP, colorFromHex(0x2ab7ca))
	drawProgressBar(screen, runnerW*2, runnerH, fallP, colorFromHex(0xfed766))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

//go:embed resources/runner.png
var runnerBytes []byte

var runnerSprite *sprite.AnimatedSprite
var idleSprite *sprite.AnimatedSprite
var fallSprite *sprite.AnimatedSprite
var runP, idleP, fallP float64

func loadSprites() error {
	img, err := imagex.FromBytes(runnerBytes)
	if err != nil {
		return err
	}
	imageGrid := imagex.NewGrid(img, runnerW, runnerH)

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
	idleSprite = sprite.NewAnimated(idleImages)
	idleSprite.SetPosition(0, 0)
	runnerSprite = sprite.NewAnimated(runImages)
	runnerSprite.SetPosition(runnerW, 0)
	fallSprite = sprite.NewAnimated(fallImages)
	fallSprite.SetPosition(runnerW*2, 0)
	return nil
}

func main() {
	game := &Game{}

	if err := loadSprites(); err != nil {
		log.Fatal(err)
	}
	animationDuration := time.Second
	go func() {
		for {
			idleP, runP, fallP = 0, 0, 0

			idleSprite.Play(animationDuration)                              // 1 iteration (default)
			runnerSprite.Play(animationDuration, animate.WithIterations(2)) // 2 iterations
			fallSprite.Play(animationDuration, animate.WithIterations(0))   // infinite iterations
			ci := animate.Float64(&idleP, 1, animationDuration)
			cr := animate.Float64(&runP, 1, animationDuration, animate.WithIterations(2))
			cf := animate.Float64(&fallP, 1, animationDuration, animate.WithIterations(0))
			time.Sleep(animationDuration * 4)
			ci()
			cr()
			cf()
			idleP, runP, fallP = 0, 0, 0

			idleSprite.Play(animationDuration)                                   // 0 seconds (default)
			runnerSprite.Play(animationDuration, animate.WithDelay(time.Second)) // 1 second
			fallSprite.Play(animationDuration, animate.WithDelay(time.Second*2)) // 2 seconds
			ci = animate.Float64(&idleP, 1, animationDuration)
			cr = animate.Float64(&runP, 1, animationDuration, animate.WithDelay(time.Second))
			cf = animate.Float64(&fallP, 1, animationDuration, animate.WithDelay(time.Second*2))
			time.Sleep(animationDuration * 3)

			idleSprite.Stop()
			runnerSprite.Stop()
			fallSprite.Stop()
			ci()
			cr()
			cf()
			idleP, runP, fallP = 0, 0, 0
			time.Sleep(animationDuration)
		}
	}()

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Animation Demo")
	if err := ebiten.RunGame(game); err != nil {
		if !errors.Is(err, ErrCloseGame) {
			log.Fatal(err)
		}
	}
}
