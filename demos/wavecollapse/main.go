package main

import (
	_ "embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/moolmanruan/ebitengine-test/grid"
	imagex "github.com/moolmanruan/ebitengine-test/image"
	"github.com/moolmanruan/ebitengine-test/sprite"
	"image"
	"image/color"
	_ "image/png"
	"log"
	"math/rand"
)

const (
	scale          = 10
	tileW, tileH   = 5, 5
	pixelW, pixelH = tileW * scale, tileH * scale
	cols, rows     = 16, 8
	screenWidth    = pixelW * cols
	screenHeight   = pixelH * rows
)

type Game struct {
	rules WaveCollapseRules
}

func generateMap(g *Game) {
	tileSprites = grid.New[*sprite.Sprite](cols, rows, nil)
	fullSet := NewSet()
	for i := 0; i < len(tileImages); i++ {
		fullSet.Add(i)
	}
	tileGrid := grid.New[Set](cols, rows, func(_, _ int) Set {
		return fullSet.Copy()
	})
	tileGrid.ForEach(func(s Set, x, y int) {
		pvs := s.Values()
		if len(pvs) == 0 {
			return
		}
		chosenTileIdx := pvs[rand.Intn(len(pvs))]

		rules := g.rules[chosenTileIdx]

		// Reduce top tile's possibilities
		topSet, err := tileGrid.Get(x, y-1)
		if err == nil {
			_ = tileGrid.Set(x, y-1, topSet.Intersection(rules.top))
		}
		// Reduce bottom tile's possibilities
		botSet, err := tileGrid.Get(x, y+1)
		if err == nil {
			_ = tileGrid.Set(x, y+1, botSet.Intersection(rules.bottom))
		}
		// Reduce right tile's possibilities
		rightSet, err := tileGrid.Get(x+1, y)
		if err == nil {
			_ = tileGrid.Set(x+1, y, rightSet.Intersection(rules.right))
		}
		// Reduce bottom tile's possibilities
		leftSet, err := tileGrid.Get(x-1, y)
		if err == nil {
			_ = tileGrid.Set(x-1, y, leftSet.Intersection(rules.left))
		}

		tile := tileImages[chosenTileIdx]
		ns := sprite.New(tile)
		ns.SetScale(scale, scale)
		ns.SetPosition(float64(x*pixelW), float64(y*pixelH))
		_ = tileSprites.Set(x, y, ns)
	})
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		generateMap(g)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	tileSprites.ForEach(func(s *sprite.Sprite, _, _ int) {
		s.Draw(screen)
	})
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

//go:embed resources/atlas.png
var imgBytes []byte

var tileImages []image.Image
var tileSprites grid.Grid[*sprite.Sprite]

func loadImages() ([]image.Image, error) {
	img, err := imagex.FromBytes(imgBytes)
	if err != nil {
		return nil, err
	}
	return imagex.TileSlice(img, tileW, tileH).List(), nil
}

var tPixels = [][2]int{{0, 0}, {4, 0}}
var bPixels = [][2]int{{0, 4}, {4, 4}}
var rPixels = [][2]int{{4, 0}, {4, 4}}
var lPixels = [][2]int{{0, 0}, {0, 4}}

func getValidTop(img image.Image, others []image.Image) Set {
	return getImagesWithMatchingPixels(img, others, tPixels, bPixels)
}
func getValidRight(img image.Image, others []image.Image) Set {
	return getImagesWithMatchingPixels(img, others, rPixels, lPixels)
}
func getValidBottom(img image.Image, others []image.Image) Set {
	return getImagesWithMatchingPixels(img, others, bPixels, tPixels)
}
func getValidLeft(img image.Image, others []image.Image) Set {
	return getImagesWithMatchingPixels(img, others, lPixels, rPixels)
}

func getImagesWithMatchingPixels(img image.Image, others []image.Image,
	imgPixels, otherPixels [][2]int) Set {
	if len(imgPixels) != len(otherPixels) {
		fmt.Println("Mismatched number of pixels")
		return nil
	}
	var pp []color.Color
	for _, v := range imgPixels {
		pp = append(pp, img.At(v[0], v[1]))
	}

	res := NewSet()
otherLoop:
	for i, other := range others {
		for oi, v := range otherPixels {
			p := other.At(v[0], v[1])
			if pp[oi] != p {
				continue otherLoop
			}
		}
		res.Add(i)
	}
	return res
}

func main() {
	var err error
	tileImages, err = loadImages()
	if err != nil {
		log.Fatal(err)
	}

	rules := make(WaveCollapseRules)
	for iIdx, img := range tileImages {
		rules[iIdx] = ValidNeighbours{
			top:    getValidTop(img, tileImages),
			right:  getValidRight(img, tileImages),
			bottom: getValidBottom(img, tileImages),
			left:   getValidLeft(img, tileImages),
		}
	}

	g := &Game{rules}
	generateMap(g) // generate an initial map

	ebiten.SetWindowSize(screenWidth*2, screenHeight*2)
	ebiten.SetWindowTitle("Wave Collapse Demo")
	if err := ebiten.RunGame(&Game{rules}); err != nil {
		log.Fatal(err)
	}
}
