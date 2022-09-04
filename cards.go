package main

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/moolmanruan/ebitengine-test/animate"
	"github.com/moolmanruan/ebitengine-test/deck"
	"github.com/moolmanruan/ebitengine-test/playingcards"
	"github.com/ungerik/go3d/float64/vec2"
	"image"
	"log"
	"math"
	"time"
)

//go:embed resources/playingcards/cards-basic.PNG
var cardsImageBytes []byte

const cardW, cardH = 17, 26

var cardsImage *ebiten.Image

func loadCardsImage() {
	img, _, err := image.Decode(bytes.NewReader(cardsImageBytes))
	if err != nil {
		log.Fatal(err)
	}
	cardsImage = ebiten.NewImageFromImage(img)
}

var suitIndex = map[playingcards.Suit]int{
	playingcards.Hearts:   0,
	playingcards.Diamonds: 1,
	playingcards.Clubs:    2,
	playingcards.Spades:   3,
}

type GameCard struct {
	playingcards.Card
	Front     *ebiten.Image
	Back      *ebiten.Image
	geom      ebiten.GeoM
	sx, sy, r float64
	pos       vec2.T
	faceDown  bool
	flipAngle float64 // in radians
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
	animate.Value(&c.pos[0], pos[0], duration, animate.WithDelay(delay))
	animate.Value(&c.pos[1], pos[1], duration, animate.WithDelay(delay))
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
	animate.Value(&c.flipAngle, dest, duration, animate.WithDelay(delay))
	return c
}

func (c *GameCard) Draw(dst *ebiten.Image) {
	face := c.Front
	if c.flipAngle > math.Pi/2 {
		face = c.Back
	}
	flipRatio := math.Cos(c.flipAngle)

	geom := c.geom
	geom.Rotate(c.r)
	geom.Scale(c.sx*flipRatio, c.sy)
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

func setupDeck() deck.Deck[*GameCard] {
	loadCardsImage()

	standardCards := playingcards.StandardDeck()
	cc := make([]*GameCard, len(standardCards))

	for i, c := range standardCards {
		cc[i] = NewGameCard(c,
			cardImage(c),
			cardBack(0))
	}
	return deck.New(cc)
}

func cardImage(card playingcards.Card) *ebiten.Image {
	x := int(card.Face) * cardW
	y := suitIndex[card.Suit] * cardH
	return cardsImage.SubImage(image.Rect(x, y, x+cardW, y+cardH)).(*ebiten.Image)
}
func cardBack(version int) *ebiten.Image {
	x := version * cardW
	y := 4 * cardH
	return cardsImage.SubImage(image.Rect(x, y, x+cardW, y+cardH)).(*ebiten.Image)
}
