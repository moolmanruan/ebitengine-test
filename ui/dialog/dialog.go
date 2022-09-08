package dialog

import (
	"github.com/hajimehoshi/ebiten/v2"
	imagex "github.com/moolmanruan/ebitengine-test/image"
	"github.com/moolmanruan/ebitengine-test/sprite"
	"image"
)

type T struct {
	sprite *sprite.Sprite
	// TODO: Handle this properly
	//tl, tc, tr *sprite.Sprite
	//ml, mc, mr *sprite.Sprite
	//bl, bc, br *sprite.Sprite
}

func New(img image.Image, top, right, bottom, left int) *T {
	ig := imagex.NewGrid3x3(img, top, right, bottom, left).List()
	return &T{
		sprite: sprite.New(ig...),
	}
}

func (d *T) Draw(dst *ebiten.Image) {
	d.sprite.Draw(dst)
}
