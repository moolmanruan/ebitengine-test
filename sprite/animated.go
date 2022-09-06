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
	s.stop = animate.Int(func(v int) {
		err := s.SetActiveImage(v)
		if err != nil {
			fmt.Println("error setting active image", err.Error())
		}
	}, 0, s.NumImages()-1, s.interval*time.Duration(s.NumImages()), animate.WithInterval(s.interval))
}

func (s *AnimatedSprite) Pause() {
	if s.stop != nil {
		s.stop()
		s.stop = nil
	}
}
