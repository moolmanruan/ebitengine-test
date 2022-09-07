package fsm

import (
	"fmt"
	"github.com/moolmanruan/ebitengine-test/fsm"
	"github.com/moolmanruan/ebitengine-test/sprite"
	"golang.org/x/exp/constraints"
	"image"
)

type FSMSprite[T constraints.Integer] struct {
	*sprite.Sprite
	sm             *fsm.FSM[T]
	stateSpriteIdx map[T]int
}

func New[T constraints.Integer]() *FSMSprite[T] {
	return &FSMSprite[T]{
		Sprite:         sprite.New(),
		sm:             fsm.New[T](),
		stateSpriteIdx: make(map[T]int, 0),
	}
}

func (s *FSMSprite[T]) AddState(state T, image image.Image) *FSMSprite[T] {
	s.sm.AddState(state)
	s.stateSpriteIdx[state] = s.AddImage(image)
	return s
}

func (s *FSMSprite[T]) AddTransition(from, to T) *FSMSprite[T] {
	s.sm.AddTransition(from, to)
	return s
}

func (s *FSMSprite[T]) To(state T) *FSMSprite[T] {
	if s.sm.To(state) {
		s.sync()
	}
	return s
}

func (s *FSMSprite[T]) SetState(state T) *FSMSprite[T] {
	s.sm.Set(state)
	s.sync()
	return s
}

func (s *FSMSprite[T]) sync() {
	if err := s.SetActiveImage(s.stateSpriteIdx[s.sm.Current()]); err != nil {
		fmt.Println("FSMSprite: Could not set the active image:", err.Error())
	}
}
