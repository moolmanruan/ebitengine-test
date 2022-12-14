package main

import (
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/moolmanruan/ebitengine-test/animate"
	"github.com/moolmanruan/ebitengine-test/deck"
	imagex "github.com/moolmanruan/ebitengine-test/image"
	"github.com/moolmanruan/ebitengine-test/playingcards"
	"github.com/ungerik/go3d/float64/vec2"
	"math"
	"time"
)

//go:embed resources/playingcards/cards-basic.PNG
var cardsImageBytes []byte

const cardW, cardH = 17, 26

var suitIndex = map[playingcards.Suit]int{
	playingcards.Hearts:   0,
	playingcards.Diamonds: 1,
	playingcards.Clubs:    2,
	playingcards.Spades:   3,
}

type GameCard struct {
	playingcards.Card
	Front      *ebiten.Image
	Back       *ebiten.Image
	geom       ebiten.GeoM
	sx, sy, r  float64
	pos        vec2.T
	faceDown   bool
	flipAngle  float64 // in radians
	cancelFlip func()
}

func NewGameCard(card playingcards.Card, frontFace *ebiten.Image, backFace *ebiten.Image) *GameCard {
	geom := ebiten.GeoM{}
	geom.Translate(-cardW/2, -cardH/2)
	return &GameCard{
		Card:  card,
		Front: frontFace,
		Back:  backFace,
		geom:  geom,
	}
}

func (c *GameCard) Position() (x, y float64) {
	return c.pos[0], c.pos[1]
}

func (c *GameCard) SetPosition(pos vec2.T, duration, delay time.Duration) *GameCard {
	animate.Float64(&c.pos[0], pos[0], duration, animate.WithDelay(delay))
	animate.Float64(&c.pos[1], pos[1], duration, animate.WithDelay(delay))
	return c
}

func (c *GameCard) SetScale(x, y float64) *GameCard {
	c.sx, c.sy = x, y
	return c
}
func (c *GameCard) SetRotation(radians float64) *GameCard {
	c.r = radians
	return c
}

func (c *GameCard) FaceUp() bool {
	return !c.faceDown
}

func (c *GameCard) SetFaceUp(up bool, duration, delay time.Duration) *GameCard {
	c.faceDown = !up
	var dest float64
	if c.faceDown {
		dest = math.Pi
	}
	if c.cancelFlip != nil {
		c.cancelFlip()
	}
	c.cancelFlip = animate.Float64(&c.flipAngle, dest, duration, animate.WithDelay(delay))
	return c
}

func (c *GameCard) Draw(dst *ebiten.Image) {
	face := c.Front
	if c.flipAngle > math.Pi/2 {
		face = c.Back
	}
	flipRatio := math.Cos(c.flipAngle) // [-1,1]

	geom := c.geom
	geom.Rotate(c.r)
	geom.Scale(c.sx*math.Abs(flipRatio), c.sy)
	x, y := c.Position()
	geom.Translate(x, y)
	dst.DrawImage(face, &ebiten.DrawImageOptions{GeoM: geom})
}

func floor(v float64) int {
	return int(math.Floor(v))
}
func ceil(v float64) int {
	return int(math.Ceil(v))
}

func (c *GameCard) In(x, y int) bool {
	chw, chh := cardW/2*c.sx, cardH/2*c.sy
	cx, cy := c.Position()
	minX, minY := floor(cx-chw), floor(cy-chh)
	maxX, maxY := ceil(cx+chw), ceil(cy+chh)
	return x >= minX && x <= maxX && y >= minY && y <= maxY
}

func setupDeck() (deck.Deck[*GameCard], error) {
	img, err := imagex.FromBytes(cardsImageBytes)
	if err != nil {
		return deck.Deck[*GameCard]{}, err
	}
	cardImgGrid := imagex.TileSlice(img, cardW, cardH)
	cardBack := cardImgGrid.At(0, 4)

	standardCards := playingcards.StandardDeck()

	cc := make([]*GameCard, len(standardCards))
	for i, c := range standardCards {
		card := cardImgGrid.At(int(c.Face), int(c.Suit)-1)
		cc[i] = NewGameCard(c,
			ebiten.NewImageFromImage(card),
			ebiten.NewImageFromImage(cardBack))
	}
	return deck.New(cc), nil
}
