package main

import (
	_ "embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/moolmanruan/ebitengine-test/deck"
	"github.com/ungerik/go3d/float64/vec2"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"image/color"
	_ "image/png"
	"log"
	"time"
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
	deck       deck.Deck[*GameCard] // the deck of cards to draw from
	deckPos    vec2.T               // the position of the deck on the screen
	cards      []*GameCard          // cards that have been drawn from the deck
	drawAmount int                  // the amount of cards to draw when the deck is clicked
}

func handleMouseClick(g *Game) {
	if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		return
	}

	x, y := ebiten.CursorPosition()
	for _, c := range g.cards {
		if c.In(x, y) {
			c.SetFaceUp(!c.FaceUp(), time.Millisecond*200, 0)
		}
	}
	if g.deck.Size() > 0 {
		if g.deck.Card(0).In(x, y) {
			drawn, newDeck := g.deck.Draw(g.drawAmount)
			g.deck = newDeck

			for ci, c := range drawn {
				i := len(g.cards)
				offset := cardSize.Scale(1.1)
				rowPos, row := i%cardsPerRow, i/cardsPerRow
				pos := vec2.T{float64(rowPos+1) * offset.X, float64(row) * offset.Y}
				moveDelay := time.Duration(ci) * time.Millisecond * 100
				moveTime := time.Second
				c.SetPosition(g.deckPos.Added(&pos), moveTime, moveDelay)
				flipTime := time.Millisecond * 200
				c.SetFaceUp(true, flipTime, moveTime+moveDelay-flipTime)
				g.cards = append(g.cards, c)
			}
		}
	}
}

func handleMouseWheel(g *Game) {
	_, w := ebiten.Wheel()
	switch {
	case w > 0:
		if g.drawAmount < 10 {
			g.drawAmount++
		}
	case w < 0:
		if g.drawAmount > 0 {
			g.drawAmount--
		}
	}
}

func handleMouse(g *Game) {
	handleMouseClick(g)
	handleMouseWheel(g)
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
	if g.deck.Size() > 0 {
		t := fmt.Sprintf("Cards left: %d  Draw rate: %d", g.deck.Size(), g.drawAmount)
		text.Draw(screen, t, mplusNormalFont, int(g.deckPos[0]-cardSize.X/2), int(g.deckPos[1]-cardSize.Y/2-5), color.White)
		g.deck.Card(0).Draw(screen)
	}
	for _, c := range g.cards {
		c.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	d := setupDeck()
	d = d.Shuffle()
	deckPos := vec2.T{50, 50}
	for _, c := range d.Cards() {
		c.SetFaceUp(false, 0, 0)
		c.SetPosition(deckPos, 0, 0)
		c.SetScale(cardScale, cardScale)
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
