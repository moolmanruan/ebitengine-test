package grid

import "errors"

type Grid[T any] struct {
	vv   []T
	w, h int
}

var ErrIndexOutOfBounds = errors.New("index out of bounds")

func New[T any](w, h int, initFn func(x, y int) T) Grid[T] {
	vv := make([]T, w*h)
	if initFn != nil {
		var i int
		for y := 0; y < h; y++ {
			for x := 0; x < w; x++ {
				vv[i] = initFn(x, y)
				i++
			}
		}
	}
	return Grid[T]{w: w, h: h, vv: vv}
}

func (g *Grid[T]) ForEach(fn func(t T, x, y int)) {
	for i, v := range g.vv {
		fn(v, i%g.w, i/g.w)
	}
}

func (g *Grid[T]) index(x, y int) (int, error) {
	if x < 0 || x >= g.w || y < 0 || y >= g.h {
		return 0, ErrIndexOutOfBounds
	}
	return x + g.w*y, nil
}

func (g *Grid[T]) Set(x, y int, value T) error {
	idx, err := g.index(x, y)
	if err != nil {
		return err
	}
	g.vv[idx] = value
	return nil
}

func (g *Grid[T]) Get(x, y int) (T, error) {
	idx, err := g.index(x, y)
	if err != nil {
		var z T
		return z, nil
	}
	return g.vv[idx], nil
}
