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

var suitIndex = map[card.Suit]int{
	card.Hearts:   0,
	card.Diamonds: 1,
	card.Clubs:    2,
	card.Spades:   3,
}

type GameCard struct {
	card.Card
	Front           *ebiten.Image
	Back            *ebiten.Image
	geom            ebiten.GeoM
	x, y, sx, sy, r float64
	faceDown        bool
	flip            float64 // Progress of card flip (0:face up -> 1:face down)
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

func (c *GameCard) SetFaceUp(up bool, duration time.Duration) *GameCard {
	c.faceDown = !up
	if duration < time.Millisecond {
		if c.faceDown {
			c.flip = 1
		} else {
			c.flip = 0
		}
		return c
	}

	if c.flip <= 0 || c.flip >= 1 {
		go animateCardFlip(c, duration)
	}
	return c
}

func animateCardFlip(c *GameCard, duration time.Duration) {
	delay := time.Millisecond * 20
	delta := float64(delay.Milliseconds()) / float64(duration.Milliseconds())

	for {
		if c.faceDown {
			c.flip += delta
			if c.flip > 1 {
				c.flip = 1
				break
			}
		} else {
			c.flip -= delta
			if c.flip < 0 {
				c.flip = 0
				break
			}
		}
		time.Sleep(delay)
	}
}

func (c *GameCard) Draw(dst *ebiten.Image) {
	face := c.Front
	if c.flip > 0.5 {
		face = c.Back
	}
	flipRatio := math.Min(math.Abs(c.flip-0.5)*2, 1)
	geom := c.geom
	geom.Rotate(c.r)
	geom.Scale(c.sx*flipRatio, c.sy)
	geom.Translate(c.x, c.y)
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
