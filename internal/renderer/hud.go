package renderer

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/neghmurken/iqube/internal/model"
)

func (r *Renderer) drawHUD(ds State) {
	x := float32(10)
	y := float32(10)
	fontSize := float32(r.style.FontSize)
	spacing := fontSize / 10

	var phase string
	switch ds.Phase {
	case model.PhasePlacement:
		phase = "PLACEMENT"
	case model.PhaseRunning:
		phase = "RUNNING"
	case model.PhaseWon:
		phase = "YOU WIN!"
	}
	rl.DrawTextEx(r.font, phase, rl.NewVector2(x, y), fontSize, spacing, rl.White)
	y += fontSize + 4

	rl.DrawTextEx(r.font, fmt.Sprintf("[1] TurnLeft: %d", ds.Inventory.TurnLeft), rl.NewVector2(x, y), fontSize, spacing, rl.White)
	y += fontSize + 2
	rl.DrawTextEx(r.font, fmt.Sprintf("[2] TurnRight: %d", ds.Inventory.TurnRight), rl.NewVector2(x, y), fontSize, spacing, rl.White)
	y += fontSize + 2
	rl.DrawTextEx(r.font, fmt.Sprintf("[3] TurnAround: %d", ds.Inventory.TurnAround), rl.NewVector2(x, y), fontSize, spacing, rl.White)
	y += fontSize + 8

	smallSize := fontSize - 4
	smallSpacing := smallSize / 10
	switch ds.Phase {
	case model.PhasePlacement:
		rl.DrawTextEx(r.font, "Space: Start  R: Reset", rl.NewVector2(x, y), smallSize, smallSpacing, rl.Gray)
	case model.PhaseRunning:
		rl.DrawTextEx(r.font, "Space: Stop", rl.NewVector2(x, y), smallSize, smallSpacing, rl.Gray)
	case model.PhaseWon:
		rl.DrawTextEx(r.font, "Space/Enter: Next level", rl.NewVector2(x, y), smallSize, smallSpacing, rl.Gray)
	}

	if r.hovered.Valid {
		label := fmt.Sprintf("%s (%d, %d)", faceNames[r.hovered.FaceIdx], r.hovered.Col, r.hovered.Row)
		dbgSize := fontSize - 4
		dbgSpacing := dbgSize / 10
		measured := rl.MeasureTextEx(r.font, label, dbgSize, dbgSpacing)
		sw := rl.GetScreenWidth()
		sh := rl.GetScreenHeight()
		rl.DrawTextEx(r.font, label, rl.NewVector2(float32(sw)-measured.X-8, float32(sh)-dbgSize-8), dbgSize, dbgSpacing, rl.Black)
	}
}
