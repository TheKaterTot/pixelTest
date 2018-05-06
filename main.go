package main

import (
	_ "image/png"
	"math/rand"
	"runtime"
	"time"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

const padding float64 = 25

var cfg = pixelgl.WindowConfig{
	Title:  "You Better Work",
	Bounds: pixel.R(0, 0, 1024, 768),
	VSync:  true,
}

type entity struct {
	Pos    pixel.Vec
	Sprite *pixel.Sprite
	Bounds pixel.Rect
	Scale  float64
}

func getInitialPos(sprite *pixel.Sprite, scale float64) pixel.Vec {
	x, y := sprite.Frame().Size().XY()
	x = x * scale
	y = y * scale
	return pixel.V(x/2+padding, y/2+padding)
}

func getBounds(sprite *entity) pixel.Rect {
	width := sprite.Sprite.Frame().W() * sprite.Scale
	height := sprite.Sprite.Frame().H() * sprite.Scale
	x := sprite.Pos.X - (width / 2.0)
	y := sprite.Pos.Y - (height / 2.0)
	return pixel.R(x, y, width+x, height+y)
}

func newEntityFromSprite(imgPath string) (*entity, error) {
	scale := 0.065
	pic, err := loadPicture(imgPath)
	if err != nil {
		return nil, err
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())
	pos := getInitialPos(sprite, scale)
	return &entity{Pos: pos, Sprite: sprite, Scale: scale}, nil
}

func placenewSprite() (*entity, error) {
	return newEntityFromSprite("./images/player.png")
}

func newMissile(pos pixel.Vec) (*entity, error) {
	imgPath := "./images/missile.png"
	scale := 0.035

	pic, err := loadPicture(imgPath)
	if err != nil {
		return nil, err
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())
	return &entity{Pos: pos, Sprite: sprite, Scale: scale}, nil
}

func newPlayerMissle(player *entity) (*entity, error) {
	return newMissile(player.Pos)
}

func playerFire(player *entity) (*entity, error) {
	return newPlayerMissle(player)
}

func isMissileOffWorld(x float64) bool {
	if x > 1024 {
		return true
	}
	return false
}

func init() {
	runtime.LockOSThread()
	rand.Seed(time.Now().Unix())
}

func main() {
	pixelgl.Run(run)
}

func run() {
	g := newGame()

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	g.running = false
	for !win.Closed() && !g.running {
		g.gameStart(win)
		if win.JustPressed(pixelgl.KeyEnter) {
			g.running = true
			break
		}
	}

	for !win.Closed() {
		for g.running && !win.Closed() {
			if !g.running {
				break
			}
			g.input(win)
			g.draw(win)
			g.update(win)
		}
		g.gameOver(win)
		if win.JustPressed(pixelgl.KeyEnter) {
			g.running = true
			g = newGame()
		}
	}
}
