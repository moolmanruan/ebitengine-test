package main

import (
	_ "embed"
	"fmt"
	"github.com/blizzy78/ebitenui"
	eimage "github.com/blizzy78/ebitenui/image"
	"github.com/blizzy78/ebitenui/widget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/moolmanruan/ebitengine-test/demos/ebitenui/resources/buttons"
	"github.com/moolmanruan/ebitengine-test/demos/ebitenui/resources/pressstart"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	_ "image/png"
	"log"
)

func main() {
	ebiten.SetWindowSize(900, 800)
	ebiten.SetWindowTitle("Ebiten UI Demo")

	ui, err := createUI()
	if err != nil {
		log.Fatal(err)
	}

	err = ebiten.RunGame(&game{
		ui: ui,
	})
	if err != nil {
		log.Print(err)
	}
}

func createUI() (*ebitenui.UI, error) {
	tt, err := opentype.Parse(pressstart.PressStart2P_ttf)
	if err != nil {
		log.Fatal(err)
	}
	const dpi = 72
	astralFF, err := opentype.NewFace(tt, &opentype.FaceOptions{Size: 24, DPI: dpi, Hinting: font.HintingFull})
	if err != nil {
		log.Fatal(err)
	}

	normalBtn9s := eimage.NewNineSlice(buttons.NormalButton(), [3]int{10, 29, 10}, [3]int{10, 29, 10})
	pressedBtn9s := eimage.NewNineSlice(buttons.PressedButton(), [3]int{10, 29, 10}, [3]int{10, 29, 10})
	btnImage := widget.ButtonImage{
		Idle:     normalBtn9s,
		Hover:    normalBtn9s,
		Pressed:  pressedBtn9s,
		Disabled: normalBtn9s,
	}

	panel9s := eimage.NewNineSlice(buttons.Panel(), [3]int{49, 2, 49}, [3]int{49, 2, 49})
	rootContainer := widget.NewContainer(
		widget.ContainerOpts.Layout(widget.NewGridLayout(
			widget.GridLayoutOpts.Columns(1),
			widget.GridLayoutOpts.Padding(widget.Insets{Top: 20, Left: 20, Right: 20, Bottom: 20}),
		)),
		widget.ContainerOpts.BackgroundImage(panel9s))

	rootContainer.AddChild(newButton("One", &btnImage, astralFF))
	rootContainer.AddChild(newButton("Two", &btnImage, astralFF))
	rootContainer.AddChild(newButton("Three", &btnImage, astralFF))

	ui := &ebitenui.UI{
		Container: rootContainer,
	}

	return ui, nil
}

func newButton(label string, img *widget.ButtonImage, ff font.Face) *widget.Button {
	return widget.NewButton(
		// specify the images to use
		widget.ButtonOpts.Image(img),

		// specify the button's text, the font face, and the color
		widget.ButtonOpts.Text(label, ff, &widget.ButtonTextColor{
			Idle: color.RGBA{0x00, 0x00, 0x00, 0xff},
		}),

		// specify that the button's text needs some padding for correct display
		widget.ButtonOpts.TextPadding(widget.Insets{
			Top:    20,
			Bottom: 20,
			Left:   30,
			Right:  30,
		}),

		widget.ButtonOpts.ClickedHandler(func(args *widget.ButtonClickedEventArgs) {
			fmt.Println(label + "CLICKED")
		}),

		// ... click handler, etc. ...
	)
}

type game struct {
	ui *ebitenui.UI
}

func (g *game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return outsideWidth, outsideHeight
}

func (g *game) Update() error {
	g.ui.Update()
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	g.ui.Draw(screen)
}
