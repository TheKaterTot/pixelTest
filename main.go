package main

import (
	"image"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

var pos pixel.Vec

func setup() {
	pos = pixel.V(50, 50)
}

func main() {
	setup()
	pixelgl.Run(run)
}

func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "You Better Work",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	pic, err := loadPicture("images/sprite-test.png")
	if err != nil {
		panic(err)
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())

	for !win.Closed() {
		win.Clear(colornames.Violet)
		sprite.Draw(win, pixel.IM.Moved(pos))

		ctrl := pixel.ZV

		if win.Pressed(pixelgl.KeyRight) && pos.X < 974 {
			ctrl.X++
		}

		if win.Pressed(pixelgl.KeyLeft) && pos.X > 50 {
			ctrl.X--
		}

		if win.Pressed(pixelgl.KeyUp) && pos.Y < 714 {
			ctrl.Y++
		}

		if win.Pressed(pixelgl.KeyDown) && pos.Y > 50 {
			ctrl.Y--
		}

		pos = ctrl.Add(pos)
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
