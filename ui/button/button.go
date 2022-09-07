package button

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/moolmanruan/ebitengine-test/sprite/fsm"
	"image"
)

type State int

const (
	Normal State = 0
	Hover  State = 1
)

type T struct {
	sprite   *fsm.FSMSprite[State]
	hoverIdx int
}

type Option func(*T)

func WithHoverImage(img image.Image) Option {
	return func(b *T) {
		b.sprite.AddState(Hover, img)
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
		sprite: fsm.New[State]().AddState(Normal, normal),
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

func (b *T) SetHover(hover bool) {
	if !hover {
		b.sprite.SetState(Normal)
		return
	}
	b.sprite.SetState(Hover)
}
