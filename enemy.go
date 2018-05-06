package main

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

func placeNewEnemy(win *pixelgl.Window) (*entity, error) {
	x, y := getCoordinates(padding+win.Bounds().W(), padding, win.Bounds().W()*2-padding, win.Bounds().H()-padding)
	enemy, err := newEnemyentityFromSprite("./images/enemy.png", x, y)
	if err != nil {
		panic(err)
	}
	return enemy, nil
}

func newEnemyentityFromSprite(imgPath string, x float64, y float64) (*entity, error) {
	pic, err := loadPicture(imgPath)
	if err != nil {
		return nil, err
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())
	pos := pixel.V(x, y)
	return &entity{Pos: pos, Sprite: sprite, Scale: 0.065}, nil
}

func isEnemyOffWorld(x float64) bool {
	if x < 0 {
		return true
	}
	return false
}
