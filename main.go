package main

import (
	"image"
	_ "image/png"
	"math/rand"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var player *Entity
var enemy *Entity

const padding float64 = 25

var cfg pixelgl.WindowConfig = pixelgl.WindowConfig{
	Title:  "You Better Work",
	Bounds: pixel.R(0, 0, 1024, 768),
	VSync:  true,
}

type Entity struct {
	Pos    pixel.Vec
	Sprite *pixel.Sprite
}

func newEntityFromSprite(imgPath string) (*Entity, error) {
	pic, err := loadPicture(imgPath)
	if err != nil {
		return nil, err
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())
	x, y := sprite.Frame().Size().XY()
	pos := pixel.V(x/2+padding, y/2+padding)
	return &Entity{Pos: pos, Sprite: sprite}, nil
}

func newEnemyEntityFromSprite(imgPath string, x float64, y float64) (*Entity, error) {
	pic, err := loadPicture(imgPath)
	if err != nil {
		return nil, err
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())
	pos := pixel.V(x, y)
	return &Entity{Pos: pos, Sprite: sprite}, nil
}

func getCoordinates(llx, lly, trx, try float64) (float64, float64) {
	a := trx - llx
	b := try - lly
	x := rand.Float64()*a + llx
	y := rand.Float64()*b + lly
	return x, y
}

func init() {
	var err error
	player, err = newEntityFromSprite("./images/sprite-test.png")
	if err != nil {
		panic(err)
	}
}

func main() {
	pixelgl.Run(run)
}

func run() {
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	var enemies []*Entity
	for i := 0; i < 4; i++ {
		x, y := getCoordinates(padding+win.Bounds().W(), padding, win.Bounds().W()*2-padding, win.Bounds().H()-padding)
		enemy, err = newEnemyEntityFromSprite("./images/enemy.png", x, y)
		if err != nil {
			panic(err)
		}
		enemies = append(enemies, enemy)
	}
	speed := 3.0
	enemySpeed := 0.8

	for !win.Closed() {
		win.Clear(colornames.Violet)
		player.Sprite.Draw(win, pixel.IM.Moved(player.Pos))

		for _, enemy := range enemies {
			enemy.Sprite.Draw(win, pixel.IM.Moved(enemy.Pos))
			enemy.Pos.X -= enemySpeed
		}

		ctrl := pixel.ZV

		if win.Pressed(pixelgl.KeyRight) && player.Pos.X < (win.Bounds().W()-padding) {
			ctrl.X += speed
		}

		if win.Pressed(pixelgl.KeyLeft) && player.Pos.X > padding {
			ctrl.X -= speed
		}

		if win.Pressed(pixelgl.KeyUp) && player.Pos.Y < (win.Bounds().H()-padding) {
			ctrl.Y += speed
		}

		if win.Pressed(pixelgl.KeyDown) && player.Pos.Y > padding {
			ctrl.Y -= speed
		}

		player.Pos = ctrl.Add(player.Pos)
		win.Update()
	}
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
