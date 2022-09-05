package button

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/moolmanruan/ebitengine-test/sprite"
	"image"
)

type T struct {
	sprite   *sprite.Sprite
	hoverIdx int
}

type Option func(*T)

func WithHoverImage(img image.Image) Option {
	return func(b *T) {
		b.sprite.AddImage(img)
		b.hoverIdx = b.sprite.NumImages() - 1
	}
}

func WithAbsolutePosition(x, y float64) Option {
	return func(b *T) {
		b.sprite.SetPosition(x, y)
	}
}

func WithSize(x, y float64) Option {
	return func(b *T) {
		b.sprite.SetSize(x, y)
	}
}

func New(normal image.Image, opts ...Option) *T {
	b := &T{
		sprite: sprite.New(normal),
	}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func (b *T) Draw(dst *ebiten.Image) {
	b.sprite.Draw(dst)
}

func (b *T) In(x, y int) bool {
	return b.sprite.In(x, y)
}

func (b *T) SetHover(hover bool) error {
	if hover {
		return b.sprite.SetActiveImage(b.hoverIdx)
	} else {
		return b.sprite.SetActiveImage(0)
	}
}
