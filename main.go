package main

import (
	"fmt"
	"image"
	_ "image/png"
	"math/rand"
	"os"
	"runtime"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

// /var player *Entity
// var enemy *Entity
// var missile *Entity
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
	Scale  float64
}

func getInitialPos(sprite *pixel.Sprite, scale float64) pixel.Vec {
	x, y := sprite.Frame().Size().XY()
	x = x * scale
	y = y * scale
	return pixel.V(x/2+padding, y/2+padding)
}

func getBounds(sprite *Entity) pixel.Rect {
	width := sprite.Sprite.Frame().W() * sprite.Scale
	height := sprite.Sprite.Frame().H() * sprite.Scale
	x := sprite.Pos.X - (width / 2.0)
	y := sprite.Pos.Y - (height / 2.0)
	return pixel.R(x, y, width+x, height+y)
}

func newEntityFromSprite(imgPath string) (*Entity, error) {
	scale := 0.065
	pic, err := loadPicture(imgPath)
	if err != nil {
		return nil, err
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())
	pos := getInitialPos(sprite, scale)
	return &Entity{Pos: pos, Sprite: sprite, Scale: scale}, nil
}

func placenewSprite() (*Entity, error) {
	return newEntityFromSprite("./images/player.png")
}

func newEnemyEntityFromSprite(imgPath string, x float64, y float64) (*Entity, error) {
	pic, err := loadPicture(imgPath)
	if err != nil {
		return nil, err
	}

	sprite := pixel.NewSprite(pic, pic.Bounds())
	pos := pixel.V(x, y)
	return &Entity{Pos: pos, Sprite: sprite, Scale: 0.065}, nil
}

func placeNewEnemy(win *pixelgl.Window) (*Entity, error) {
	x, y := getCoordinates(padding+win.Bounds().W(), padding, win.Bounds().W()*2-padding, win.Bounds().H()-padding)
	enemy, err := newEnemyEntityFromSprite("./images/enemy.png", x, y)
	if err != nil {
		panic(err)
	}
	return enemy, nil
}

func makeEnemies(g *game, win *pixelgl.Window, number int) {
	for i := 0; i < number; i++ {
		enemy, err := placeNewEnemy(win)
		g.enemies = append(g.enemies, enemy)
		if err != nil {
			panic(err)
		}
	}
}

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

func isEnemyOffWorld(x float64) bool {
	if x < 0 {
		return true
	}
	return false
}

func filterDeadEnemies(enemies []*Entity) []*Entity {
	for i, enemy := range enemies {
		if isEnemyOffWorld(enemy.Pos.X) {
			enemies = append(enemies[:i], enemies[i+1:]...)
		}
	}
	return enemies
}

func overlap(sprite *Entity, sprite2 *Entity) bool {
	intersection := getBounds(sprite).Intersect(getBounds(sprite2))
	if intersection.W() == 0 && intersection.H() == 0 {
		return false
	}
	return true
}

func init() {
	runtime.LockOSThread()
	running = true
}

type game struct {
	score    int64
	player   *Entity
	enemies  []*Entity
	missiles []*Entity
}

func newGame() *game {
	player, _ := placenewSprite()

	return &game{
		int64(0),
		player,
		[]*Entity{},
		[]*Entity{},
	}
}

func (g *game) input(win *pixelgl.Window) {
	win.SetClosed(win.JustPressed(pixelgl.KeyEscape))

	speed := 3.0
	ctrl := pixel.ZV

	if win.Pressed(pixelgl.KeyRight) && g.player.Pos.X < (win.Bounds().W()-padding) {
		ctrl.X += speed
	}
	if win.Pressed(pixelgl.KeyLeft) && g.player.Pos.X > padding {
		ctrl.X -= speed
	}

	if win.Pressed(pixelgl.KeyUp) && g.player.Pos.Y < (win.Bounds().H()-padding) {
		ctrl.Y += speed
	}

	if win.Pressed(pixelgl.KeyDown) && g.player.Pos.Y > padding {
		ctrl.Y -= speed
	}

	g.player.Pos = ctrl.Add(g.player.Pos)

}

func (g *game) draw(win *pixelgl.Window) {
	win.Clear(colornames.Cornflowerblue)
	g.player.Sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, g.player.Scale).Moved(g.player.Pos))
	for _, enemy := range g.enemies {
		enemy.Sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, enemy.Scale).Moved(enemy.Pos))
	}
}

func (g *game) update(win *pixelgl.Window) {
	for _, enemy := range g.enemies {
		if overlap(g.player, enemy) {
			running = false
			break
		}
	}
	g.enemies = filterDeadEnemies(g.enemies)
	makeEnemies(g, win, 4-len(g.enemies))
	enemySpeed := 1.5
	for _, enemy := range g.enemies {
		enemy.Pos.X -= enemySpeed
	}

	txtvec := getTextCoordinates(win)

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(txtvec, basicAtlas)
	basicTxt.Color = colornames.Black
	fmt.Fprintf(basicTxt, "Score: %d", g.score)
	basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 2))

	win.Update()
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

	for !win.Closed() {
		for running && !win.Closed() {
			if !running {
				break
			}
			g.input(win)
			g.draw(win)
			g.update(win)
		}

		win.Clear(colornames.Black)
		basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
		basicTxt := text.New(pixel.V(100, 500), basicAtlas)
		fmt.Fprintln(basicTxt, "Press Enter to Start")
		basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 4))
		win.Update()
		if win.JustPressed(pixelgl.KeyEnter) {
			running = true
			g = newGame()
		}
		// }
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
