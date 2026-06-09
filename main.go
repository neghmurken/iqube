package main

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/neghmurken/iqube/internal/game"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

func main() {
	rl.InitWindow(screenWidth, screenHeight, "iQube")
	defer rl.CloseWindow()

	rl.SetTargetFPS(60)

	g := game.New()

	for !rl.WindowShouldClose() {
		g.Update()

		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		g.Draw()
		rl.EndDrawing()
	}
}
