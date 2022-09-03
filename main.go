package main

import (
	_ "embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/moolmanruan/ebitengine-test/deck"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	_ "image/png"
	"log"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

type Point2D struct {
	X, Y float64
}

func (p Point2D) Scale(v float64) Point2D {
	return Point2D{X: p.X * v, Y: p.Y * v}
}

type Game struct {
	mouseState MouseState           // the mouse state the last time it changed
	deck       deck.Deck[*GameCard] // the deck of cards to draw from
	deckPos    Point2D              // the position of the deck on the screen
	cards      []*GameCard          // cards that have been drawn from the deck
	drawAmount int                  // the amount of cards to draw when the deck is clicked
}

type MouseState struct {
	position    int
	leftPressed bool
	wheelUpDown float64
}

func handleMouseStateChange(g *Game) {
	if g.mouseState.leftPressed {
		x, y := ebiten.CursorPosition()
		for _, c := range g.cards {
			if c.In(x, y) {
				c.FaceDown = !c.FaceDown
			}
		}
		if g.deck.Size() > 0 {
			if g.deck.Card(0).In(x, y) {
				drawn, newDeck := g.deck.Draw(g.drawAmount)
				g.deck = newDeck
				g.cards = append(g.cards, drawn...)
			}
		}
	}
	if g.mouseState.wheelUpDown > 0 {
		if g.drawAmount < 10 {
			g.drawAmount++
		}
	} else if g.mouseState.wheelUpDown < 0 {
		if g.drawAmount > 0 {
			g.drawAmount--
		}
	}
}

func handleMouse(g *Game) {
	_, wy := ebiten.Wheel()
	ms := MouseState{
		leftPressed: ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft),
		wheelUpDown: wy,
	}
	if ms != g.mouseState {
		g.mouseState = ms
		handleMouseStateChange(g)
	}
}

var (
	mplusNormalFont font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    14,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Update() error {
	handleMouse(g)
	return nil
}

const cardScale = 2
const cardsPerRow = 13

var cardSize = Point2D{cardW, cardH}.Scale(cardScale)

func (g *Game) Draw(screen *ebiten.Image) {
	p := g.deckPos
	if g.deck.Size() > 0 {
		t := fmt.Sprintf("Cards left: %d  Draw rate: %d", g.deck.Size(), g.drawAmount)
		text.Draw(screen, t, mplusNormalFont, int(p.X-cardSize.X/2), int(p.Y-cardSize.Y/2), color.White)
		c := g.deck.Card(0).SetPosition(p.X, p.Y)
		screen.DrawImage(c.Back, c.DrawOptions())
	}
	for i, c := range g.cards {
		offset := cardSize.Scale(1.1)
		pos, row := i%cardsPerRow, i/cardsPerRow
		card := c.SetPosition(p.X+float64(pos+1)*offset.X, p.Y+float64(row)*offset.Y)
		var face *ebiten.Image
		if card.FaceDown {
			face = card.Back
		} else {
			face = card.Front
		}
		screen.DrawImage(face, card.DrawOptions())
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	d := setupDeck()
	d = d.Shuffle()
	deckPos := Point2D{50, 50}
	for _, c := range d.Cards() {
		c.SetPosition(deckPos.X, deckPos.Y).SetScale(cardScale)
	}
	game := &Game{
		deck:       d,
		deckPos:    deckPos,
		drawAmount: 2,
	}

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Card draw")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
