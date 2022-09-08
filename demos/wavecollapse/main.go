package main

import (
	_ "embed"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	cols, rows     = 10, 10
	screenWidth    = pixelW * cols
	screenHeight   = pixelH * rows
)

type Game struct {
	rules WaveCollapseRules
}

func generateMap(g *Game) {
	sprites = make([]*sprite.Sprite, cols*rows)
	fullSet := NewSet()
	for i := 0; i < len(tileImages); i++ {
		fullSet.Add(i)
	}
	spritesPossibleIdx := make([]Set, cols*rows)
	for i := 0; i < len(spritesPossibleIdx); i++ {
		spritesPossibleIdx[i] = fullSet.Copy()
	}

	indexFor := func(x, y int) int {
		if x < 0 || x >= cols || y < 0 || y >= rows {
			return -1
		}
		return x + y*cols
	}
	for si := range sprites {
		possibleIndexes := spritesPossibleIdx[si].Values()
		if len(possibleIndexes) == 0 {
			return
		}
		chosenTileIdx := possibleIndexes[rand.Intn(len(possibleIndexes))]

		rules := g.rules[chosenTileIdx]

		sx, sy := si%cols, si/cols
		tI, bI, rI, lI := indexFor(sx, sy-1), indexFor(sx, sy+1), indexFor(sx+1, sy), indexFor(sx-1, sy)
		// Reduce top tile's possibilities
		if tI != -1 {
			spritesPossibleIdx[tI] = spritesPossibleIdx[tI].Intersection(rules.top)
		}
		// Reduce left tile's possibilities
		if lI != -1 {
			spritesPossibleIdx[lI] = spritesPossibleIdx[lI].Intersection(rules.left)
		}
		// Reduce right tile's possibilities
		if rI != -1 {
			spritesPossibleIdx[rI] = spritesPossibleIdx[rI].Intersection(rules.right)
		}
		// Reduce bottom tile's possibilities
		if bI != -1 {
			spritesPossibleIdx[bI] = spritesPossibleIdx[bI].Intersection(rules.bottom)
		}

		tile := tileImages[chosenTileIdx]
		ns := sprite.New(tile)
		ns.SetScale(scale, scale)
		ns.SetPosition(float64(si%cols)*pixelW, float64(si/cols)*pixelH)
		sprites[si] = ns
	}
}

func (g *Game) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		generateMap(g)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for _, s := range sprites {
		s.Draw(screen)
	}
}

func (g *Game) Layout(_, _ int) (int, int) {
	return screenWidth, screenHeight
}

//go:embed resources/atlas.png
var imgBytes []byte

var tileImages []image.Image
var sprites []*sprite.Sprite

func loadImages() ([]image.Image, error) {
	img, err := imagex.FromBytes(imgBytes)
	if err != nil {
		return nil, err
	}
	return imagex.NewGrid(img, tileW, tileH).List(), nil
}

var tPixels = [][2]int{{0, 0}, {2, 0}, {4, 0}}
var bPixels = [][2]int{{0, 4}, {2, 4}, {4, 4}}
var rPixels = [][2]int{{4, 0}, {4, 2}, {4, 4}}
var lPixels = [][2]int{{0, 0}, {0, 2}, {0, 4}}

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
