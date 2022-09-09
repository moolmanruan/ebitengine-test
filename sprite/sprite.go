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

func New(imgs ...image.Image) *Sprite {
	s := &Sprite{
		sx: 1, sy: 1,
	}
	s.recalculateDrawOpts()
	for _, img := range imgs {
		s.AddImage(img)
	}
	return s
}

func (s *Sprite) Draw(dst *ebiten.Image) {
	if s == nil {
		return
	}
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

func (s *Sprite) AddImage(img image.Image) int {
	s.images = append(s.images, ebiten.NewImageFromImage(img))
	return len(s.images) - 1
}

func (s *Sprite) NumImages() int {
	return len(s.images)
}

func (s *Sprite) SetActiveImage(index int) error {
	if index < 0 || index >= len(s.images) {
		return errors.New("index out of bounds")
	}
	s.activeImg = index
	return nil
}

func (s *Sprite) ActiveImage() int {
	return s.activeImg
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

func (s *Sprite) GetSize() (x, y float64) {
	w, h := s.images[0].Size()
	return float64(w) * s.sx, float64(h) * s.sy
}

// In indicates whether the given screen coordinates are in the Sprite's rectangle.
func (s *Sprite) In(x, y int) bool {
	xf, yf := float64(x), float64(y)
	w, h := s.images[0].Size()
	xMin, xMax := s.px, s.px+float64(w)*s.sx
	yMin, yMax := s.py, s.py+float64(h)*s.sy
	return xf >= xMin && xf <= xMax && yf >= yMin && yf <= yMax
}
