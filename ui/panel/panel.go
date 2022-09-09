package panel

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/moolmanruan/ebitengine-test/grid"
	imagex "github.com/moolmanruan/ebitengine-test/image"
	"github.com/moolmanruan/ebitengine-test/sprite"
	"image"
)

type T struct {
	spriteGrid grid.Grid[*sprite.Sprite]
	rect       image.Rectangle
	horDist    [3]float64
	verDist    [3]float64
}

func New(imgGrid imagex.Grid, rectangle image.Rectangle) *T {
	var hd, vd [3]float64
	sg := grid.New(3, 3, func(x, y int) *sprite.Sprite {
		img := imgGrid.ImageAt(x, y)
		b := img.Bounds()
		hd[x] = float64(b.Dx())
		vd[x] = float64(b.Dy())
		return sprite.New(img)
	})

	return &T{
		spriteGrid: sg,
		horDist:    hd,
		verDist:    vd,
		rect:       rectangle,
	}
}

func (d *T) SetRect(rect image.Rectangle) {
	d.rect = rect
}

func (d *T) Draw(dst *ebiten.Image) {
	b := d.rect.Bounds()
	minX, minY := float64(b.Min.X), float64(b.Min.Y)
	maxX, maxY := float64(b.Max.X), float64(b.Max.Y)
	xx := [4]float64{minX, minX + d.horDist[0], maxX - d.horDist[2], maxX}
	yy := [4]float64{minY, minY + d.verDist[0], maxY - d.verDist[2], maxY}
	d.spriteGrid.ForEach(func(s *sprite.Sprite, x, y int) {
		s.SetPosition(xx[x], yy[y])
		s.SetSize(xx[x+1]-xx[x], yy[y+1]-yy[y])
		s.Draw(dst)
	})
}
