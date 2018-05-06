package main

import (
	"fmt"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
	"golang.org/x/image/font/basicfont"
)

type gameState struct {
	enemies  []*entity
	missiles []*entity
}

type game struct {
	score   int64
	player  *entity
	current gameState
	next    gameState
}

func newGameState() gameState {
	return gameState{
		[]*entity{},
		[]*entity{},
	}
}

func newGame() *game {
	player, _ := placenewSprite()

	return &game{
		score:   int64(0),
		player:  player,
		current: newGameState(),
		next:    newGameState(),
	}
}

func (g *game) swapStates() {
	g.current = g.next
	g.next = newGameState()
}

func (g *game) makeEnemies(win *pixelgl.Window, max int) {
	number := max - len(g.current.enemies)

	for i := 0; i < number; i++ {
		enemy, err := placeNewEnemy(win)
		if err != nil {
			panic(err)
		}
		g.current.enemies = append(g.current.enemies, enemy)
	}
}

func (g *game) updateEnemies() {
	enemySpeed := 1.5

	for _, enemy := range g.current.enemies {
		enemy.Pos.X -= enemySpeed
	}
}

func (g *game) updateMissiles() {
	missileSpeed := 2.5

	for _, missile := range g.current.missiles {
		missile.Pos.X += missileSpeed
	}
}

func (g *game) gameOver(win *pixelgl.Window) {
	win.Clear(colornames.Black)
	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(pixel.V(100, 500), basicAtlas)
	fmt.Fprintln(basicTxt, "GAME OVER")
	fmt.Fprintln(basicTxt, "Press Enter to Start")
	basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 4))
	win.Update()
}

func (g *game) displayScore(win *pixelgl.Window) {
	txtvec := getTextCoordinates(win)

	basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
	basicTxt := text.New(txtvec, basicAtlas)
	basicTxt.Color = colornames.Black
	fmt.Fprintf(basicTxt, "Score: %d", g.score)
	basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 2))
}

func (g *game) checkPlayer() {
	for _, enemy := range g.current.enemies {
		if overlap(g.player, enemy) {
			running = false
			break
		}
	}
}

func (g *game) filterDeadEnemies() ([]*entity, int64) {
	enemies := []*entity{}
	hits := int64(0)
	for _, enemy := range g.current.enemies {
		if !isEnemyOffWorld(enemy.Pos.X) && !anyOverlap(enemy, g.current.missiles) {
			enemies = append(enemies, enemy)
		}
		if anyOverlap(enemy, g.current.missiles) {
			hits++
		}
	}
	return enemies, hits
}

func (g *game) filterDeadMissiles() []*entity {
	missiles := []*entity{}
	for _, missile := range g.current.missiles {
		if !isMissileOffWorld(missile.Pos.X) && !anyOverlap(missile, g.current.enemies) {
			missiles = append(missiles, missile)
		}
	}
	return missiles
}

func (g *game) draw(win *pixelgl.Window) {
	win.Clear(colornames.Cornflowerblue)
	g.player.Sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, g.player.Scale).Moved(g.player.Pos))
	for _, enemy := range g.current.enemies {
		enemy.Sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, enemy.Scale).Moved(enemy.Pos))
	}
	for _, missile := range g.current.missiles {
		missile.Sprite.Draw(win, pixel.IM.Scaled(pixel.ZV, missile.Scale).Moved(missile.Pos))
	}
}

func (g *game) update(win *pixelgl.Window) {
	g.checkPlayer()
	var hits int64
	g.next.enemies, hits = g.filterDeadEnemies()
	g.next.missiles = g.filterDeadMissiles()
	g.score += hits

	g.swapStates()
	g.makeEnemies(win, 4)
	g.updateEnemies()
	g.updateMissiles()

	g.displayScore(win)
	win.Update()
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

	if win.JustPressed(pixelgl.KeySpace) {
		missile, _ := playerFire(g.player)
		g.current.missiles = append(g.current.missiles, missile)
	}

}
