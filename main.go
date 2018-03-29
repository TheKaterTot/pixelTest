package main

import (
	"image"
	_ "image/png"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
)

func main() {
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
	v := pixel.V(50, 50)

	for !win.Closed() {
		win.Clear(colornames.Violet)
		sprite.Draw(win, pixel.IM.Moved(v))

		ctrl := pixel.ZV

		if win.Pressed(pixelgl.KeyRight) {
			ctrl.X++
		}

		if win.Pressed(pixelgl.KeyLeft) {
			ctrl.X--
		}

		if win.Pressed(pixelgl.KeyUp) {
			ctrl.Y++
		}

		if win.Pressed(pixelgl.KeyDown) {
			ctrl.Y--
		}

		v = ctrl.Add(v)
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
