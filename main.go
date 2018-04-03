package main

import (
	"image"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var player *Entity

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

	speed := 3.0

	for !win.Closed() {
		win.Clear(colornames.Violet)
		player.Sprite.Draw(win, pixel.IM.Moved(player.Pos))

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
