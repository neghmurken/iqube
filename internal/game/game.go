package game

import rl "github.com/gen2brain/raylib-go/raylib"

type Game struct{}

func New() *Game {
	return &Game{}
}

func (g *Game) Update() {}

func (g *Game) Draw() {
	rl.DrawText("iQube", 10, 10, 20, rl.White)
}
