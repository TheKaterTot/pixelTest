package main

import (
	"fmt"
	"image"
	_ "image/png"
	"math/rand"
	"os"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

var player *Entity
var enemy *Entity
var running bool

const padding float64 = 25

var cfg pixelgl.WindowConfig = pixelgl.WindowConfig{
	Title:  "You Better Work",
	Bounds: pixel.R(0, 0, 1024, 768),
	VSync:  true,
}

type Entity struct {
	Pos    pixel.Vec
	Sprite *pixel.Sprite
	Bounds pixel.Rect
}

func getInitialPos(sprite *pixel.Sprite) pixel.Vec {
	x, y := sprite.Frame().Size().XY()
	return pixel.V(x/2+padding, y/2+padding)
}

func getBounds(sprite *Entity) pixel.Rect {
	return pixel.R(
		sprite.Pos.X,
		sprite.Pos.Y,
		sprite.Pos.X+sprite.Sprite.Frame().W(),
		sprite.Pos.Y+sprite.Sprite.Frame().H())
}

func newEntityFromSprite(imgPath string) (*Entity, error) {
	pic, err := loadPicture(imgPath)
	if err != nil {
		return nil, err
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())
	pos := getInitialPos(sprite)
	return &Entity{Pos: pos, Sprite: sprite}, nil
}

func placenewSprite() {
	var err error
	player, err = newEntityFromSprite("./images/sprite-test.png")
	if err != nil {
		panic(err)
	}
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

func placeNewEnemy(win *pixelgl.Window) (*Entity, error) {
	x, y := getCoordinates(padding+win.Bounds().W(), padding, win.Bounds().W()*2-padding, win.Bounds().H()-padding)
	enemy, err := newEnemyEntityFromSprite("./images/enemy.png", x, y)
	if err != nil {
		panic(err)
	}
	return enemy, nil
}

func getCoordinates(llx, lly, trx, try float64) (float64, float64) {
	a := trx - llx
	b := try - lly
	x := rand.Float64()*a + llx
	y := rand.Float64()*b + lly
	return x, y
}

func isOffWorld(x float64) bool {
	if x < 0 {
		return true
	}
	return false
}

func filterDeadEnemies(enemies []*Entity) []*Entity {
	var newList []*Entity
	for _, enemy := range enemies {
		if !isOffWorld(enemy.Pos.X) && !overlap(player, enemy) {
			newList = append(newList, enemy)
		}
	}
	return newList
}

func overlap(sprite *Entity, sprite2 *Entity) bool {
	sprite.Bounds = getBounds(sprite)
	sprite2.Bounds = getBounds(sprite2)
	intersection := sprite.Bounds.Intersect(sprite2.Bounds)
	if intersection.W() == 0 && intersection.H() == 0 {
		return false
	}
	return true
}

func init() {
	running = true
	placenewSprite()
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
		enemy, err := placeNewEnemy(win)
		enemies = append(enemies, enemy)
		if err != nil {
			panic(err)
		}
	}

	for !win.Closed() {
		for running {
			for _, enemy := range enemies {
				if overlap(player, enemy) {
					running = false
					break
				}
			}
			if !running {
				break
			}
			update(win, player, enemies)
			win.Clear(colornames.Violet)
			player.Sprite.Draw(win, pixel.IM.Moved(player.Pos))
			for _, enemy := range enemies {
				enemy.Sprite.Draw(win, pixel.IM.Moved(enemy.Pos))
			}
			win.Update()
		}

		win.Clear(colornames.Black)
		basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		basicTxt := text.New(pixel.V(100, 500), basicAtlas)
		fmt.Fprintln(basicTxt, "Press Space to Start")
		basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 4))
		win.Update()
		if win.JustPressed(pixelgl.KeySpace) {
			running = true
			placenewSprite()
		}
	}
}

func update(win *pixelgl.Window, player *Entity, enemies []*Entity) {
	speed := 3.0
	enemySpeed := 1.0
	enemies = filterDeadEnemies(enemies)
	if len(enemies) < 4 {
		newEnemy, _ := placeNewEnemy(win)
		enemies = append(enemies, newEnemy)
	}

	for _, enemy := range enemies {
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
