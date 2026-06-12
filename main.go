package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/neghmurken/iqube/internal/game"
)

const (
	screenWidth  = 1024
	screenHeight = 768
)

func main() {
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(screenWidth, screenHeight, "iQube")
	defer rl.CloseWindow()

	rl.SetTargetFPS(120)

	g := game.New()
	defer g.Close()

	for !rl.WindowShouldClose() {
		g.Update()

		rl.BeginDrawing()
		g.Draw()
		rl.EndDrawing()
	}
}
