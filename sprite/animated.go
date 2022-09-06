package sprite

import (
	"fmt"
	"github.com/moolmanruan/ebitengine-test/animate"
	"image"
	"time"
)

type AnimatedSprite struct {
	*Sprite
	stop     func()
	interval time.Duration
}

func NewAnimated(imgs []image.Image, interval time.Duration) *AnimatedSprite {
	return &AnimatedSprite{
		Sprite:   New(imgs...),
		interval: interval,
	}
}

func (s *AnimatedSprite) Play() {
	s.Pause()
	s.stop = animate.Int(func(newImage int) {
		err := s.SetActiveImage(newImage)
		if err != nil {
			fmt.Println("error setting active image", err.Error())
		}
	}, 0, s.NumImages()-1, animate.WithInterval(s.interval))
}

func (s *AnimatedSprite) Pause() {
	if s.stop != nil {
		s.stop()
		s.stop = nil
	}
}
