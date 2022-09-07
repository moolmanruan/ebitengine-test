package fsm

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/moolmanruan/ebitengine-test/fsm"
	"github.com/moolmanruan/ebitengine-test/sprite"
	"golang.org/x/exp/constraints"
	"image"
)

type FSMSprite[T constraints.Integer] struct {
	sm      *fsm.StateMachine[T]
	sprites map[T]*sprite.Sprite
}

func New[T constraints.Integer]() *FSMSprite[T] {
	return &FSMSprite[T]{
		sm:      fsm.New[T](),
		sprites: make(map[T]*sprite.Sprite, 0),
	}
}

func (s *FSMSprite[T]) AddState(state T, image image.Image) *FSMSprite[T] {
	s.sm.AddState(state)
	s.sprites[state] = sprite.New(image)
	return s
}

func (s *FSMSprite[T]) SetState(state T) *FSMSprite[T] {
	s.sm.Set(state)
	return s
}

func (s *FSMSprite[T]) Draw(dst *ebiten.Image) {
	if sp, ok := s.sprites[s.sm.Current()]; ok {
		sp.Draw(dst)
	}
}

func (s *FSMSprite[T]) In(x, y int) bool {
	for _, img := range s.sprites {
		return img.In(x, y)
	}
	return false
}

func (s *FSMSprite[T]) SetPosition(x, y float64) {
	for _, img := range s.sprites {
		img.SetPosition(x, y)
	}
}

func (s *FSMSprite[T]) SetSize(x, y float64) {
	for _, img := range s.sprites {
		img.SetSize(x, y)
	}
}
