package main

import (
	"image"
	"math/rand"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func getCoordinates(llx, lly, trx, try float64) (float64, float64) {
	a := trx - llx
	b := try - lly
	x := rand.Float64()*a + llx
	y := rand.Float64()*b + lly
	return x, y
}

func getTextCoordinates(win *pixelgl.Window) pixel.Vec {
	x := padding
	y := win.Bounds().H() - padding
	return pixel.V(x, y)
}

func overlap(sprite *entity, sprite2 *entity) bool {
	intersection := getBounds(sprite).Intersect(getBounds(sprite2))
	if intersection.W() == 0 && intersection.H() == 0 {
		return false
	}
	return true
}

func anyOverlap(entity *entity, others []*entity) bool {
	for _, other := range others {
		if overlap(other, entity) {
			return true
		}
	}
	return false
}

func loadPicture(path string) (pixel.Picture, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return pixel.PictureDataFromImage(img), nil
}
