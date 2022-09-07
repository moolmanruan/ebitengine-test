package sprite

import (
	"fmt"
	"github.com/moolmanruan/ebitengine-test/animate"
	"image"
	"time"
)

type AnimatedSprite struct {
	*Sprite
	stop func()
}

func NewAnimated(imgs []image.Image) *AnimatedSprite {
	return &AnimatedSprite{
		Sprite: New(imgs...),
	}
}
func (s *AnimatedSprite) Playing() bool {
	return s.stop != nil
}

func (s *AnimatedSprite) Play(duration time.Duration, opts ...animate.Option) {
	s.Stop()
	s.stop = animate.Int(func(newImage int) {
		err := s.SetActiveImage(newImage)
		if err != nil {
			fmt.Println("error setting active image", err.Error())
		}
	}, 0, s.NumImages()-1, duration, opts...)
}

func (s *AnimatedSprite) Stop() {
	if s.stop != nil {
		s.stop()
		s.stop = nil
	}
}
