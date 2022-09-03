package main

import (
	"bytes"
	_ "embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/moolmanruan/ebitengine-test/card"
	"github.com/moolmanruan/ebitengine-test/deck"
	"image"
	"log"
	"math"
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

var suitIndex = map[card.Suit]int{
	card.Hearts:   0,
	card.Diamonds: 1,
	card.Clubs:    2,
	card.Spades:   3,
}

type GameCard struct {
	card.Card
	Front      *ebiten.Image
	Back       *ebiten.Image
	FaceDown   bool
	geom       ebiten.GeoM
	x, y, s, r float64
}

func NewGameCard(card card.Card, frontFace *ebiten.Image, backFace *ebiten.Image) *GameCard {
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
	return c.x, c.y
}
func (c *GameCard) SetPosition(x, y float64) *GameCard {
	c.x, c.y = x, y
	return c
}
func (c *GameCard) SetScale(scale float64) *GameCard {
	c.s = scale
	return c
}
func (c *GameCard) SetRotation(radians float64) *GameCard {
	c.r = radians
	return c
}

func (c *GameCard) DrawOptions() *ebiten.DrawImageOptions {
	geom := c.geom
	geom.Rotate(c.r)
	geom.Scale(c.s, c.s)
	geom.Translate(c.x, c.y)
	return &ebiten.DrawImageOptions{GeoM: geom}
}

func floor(v float64) int {
	return int(math.Floor(v))
}
func ceil(v float64) int {
	return int(math.Ceil(v))
}

func (c *GameCard) In(x, y int) bool {
	chw, chh := cardW/2*c.s, cardH/2*c.s
	minX, minY := floor(c.x-chw), floor(c.y-chh)
	maxX, maxY := ceil(c.x+chw), ceil(c.y+chh)
	return x >= minX && x <= maxX && y >= minY && y <= maxY
}

func setupDeck() deck.Deck[*GameCard] {
	loadCardsImage()

	standardCards := card.StandardDeck()
	cc := make([]*GameCard, len(standardCards))

	for i, c := range standardCards {
		cc[i] = NewGameCard(c,
			cardImage(c),
			cardBack(0))
	}
	return deck.New(cc)
}

func cardImage(card card.Card) *ebiten.Image {
	x := int(card.Face) * cardW
	y := suitIndex[card.Suit] * cardH
	return cardsImage.SubImage(image.Rect(x, y, x+cardW, y+cardH)).(*ebiten.Image)
}
func cardBack(version int) *ebiten.Image {
	x := version * cardW
	y := 4 * cardH
	return cardsImage.SubImage(image.Rect(x, y, x+cardW, y+cardH)).(*ebiten.Image)
}
