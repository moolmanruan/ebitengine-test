package sprite

import (
	"bytes"
	"errors"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type Sprite struct {
	images    []*ebiten.Image
	activeImg int
	drawOpts  *ebiten.DrawImageOptions
	px, py    float64
	sx, sy    float64
}

func NewFromBytes(imgBytes []byte) (*Sprite, error) {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return nil, err
	}
	return New(img), nil
}

func New(img image.Image) *Sprite {
	s := new(Sprite)
	s.AddImage(img)
	s.recalculateDrawOpts()
	return s
}

func (s *Sprite) Draw(dst *ebiten.Image) {
	dst.DrawImage(s.images[s.activeImg], s.drawOpts)
}

func (s *Sprite) recalculateDrawOpts() {
	g := ebiten.GeoM{}
	g.Scale(s.sx, s.sy)
	g.Translate(s.px, s.py)
	s.drawOpts = &ebiten.DrawImageOptions{GeoM: g}
}

func (s *Sprite) AddImageFromBytes(imgBytes []byte) error {
	img, _, err := image.Decode(bytes.NewReader(imgBytes))
	if err != nil {
		return err
	}
	s.AddImage(img)
	return nil
}

func (s *Sprite) AddImage(img image.Image) {
	s.images = append(s.images, ebiten.NewImageFromImage(img))
}
func (s *Sprite) SetActiveImage(index int) error {
	if index < 0 || index >= len(s.images) {
		return errors.New("index out of bounds")
	}
	s.activeImg = index
	return nil
}

func (s *Sprite) SetPosition(x, y float64) {
	s.px, s.py = x, y
	s.recalculateDrawOpts()
}

func (s *Sprite) SetScale(x, y float64) {
	s.sx, s.sy = x, y
	s.recalculateDrawOpts()
}

func (s *Sprite) SetSize(x, y float64) {
	w, h := s.images[0].Size()
	s.SetScale(x/float64(w), y/float64(h))
}

// In indicates whether the given screen coordinates are in the Sprite's rectangle.
func (s *Sprite) In(x, y int) bool {
	xf, yf := float64(x), float64(y)
	w, h := s.images[0].Size()
	xMin, xMax := s.px, s.px+float64(w)*s.sx
	yMin, yMax := s.py, s.py+float64(h)*s.sy
	return xf >= xMin && xf <= xMax && yf >= yMin && yf <= yMax
}
