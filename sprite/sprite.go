package sprite

import (
	"bytes"
	"github.com/hajimehoshi/ebiten/v2"
	"image"
)

type Sprite struct {
	img      *ebiten.Image
	drawOpts *ebiten.DrawImageOptions
	px, py   float64
	sx, sy   float64
}

func NewFromBytes(bb []byte) (*Sprite, error) {
	img, _, err := image.Decode(bytes.NewReader(bb))
	if err != nil {
		return nil, err
	}
	return NewFromImage(img), nil
}

func NewFromImage(img image.Image) *Sprite {
	s := &Sprite{
		img: ebiten.NewImageFromImage(img),
	}
	s.recalculateDrawOpts()
	return s
}

func (s *Sprite) Draw(dst *ebiten.Image) {
	dst.DrawImage(s.img, s.drawOpts)
}

func (s *Sprite) recalculateDrawOpts() {
	g := ebiten.GeoM{}
	g.Scale(s.sx, s.sy)
	g.Translate(s.px, s.py)
	s.drawOpts = &ebiten.DrawImageOptions{GeoM: g}
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
	w, h := s.img.Size()
	s.SetScale(x/float64(w), y/float64(h))
}

// In indicates whether the given screen coordinates are in the Sprite's rectangle.
func (s *Sprite) In(x, y int) bool {
	xf, yf := float64(x), float64(y)
	w, h := s.img.Size()
	xMin, xMax := s.px, s.px+float64(w)*s.sx
	yMin, yMax := s.py, s.py+float64(h)*s.sy
	return xf >= xMin && xf <= xMax && yf >= yMin && yf <= yMax
}
