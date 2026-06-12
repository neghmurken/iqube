package game

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/neghmurken/iqube/internal/level"
	"github.com/neghmurken/iqube/internal/model"
	"github.com/neghmurken/iqube/internal/renderer"
	"github.com/neghmurken/iqube/internal/simulation"
)

type Game struct {
	cube       *model.Cube
	renderer   *renderer.Renderer
	sim        *simulation.Simulation
	levels     []model.Level
	levelIndex int
}

func New() *Game {
	style, err := renderer.LoadStyle("assets/theme.json")
	if err != nil {
		style = renderer.DefaultStyle()
	}

	levels, err := level.LoadAll("assets/levels")
	if err != nil || len(levels) == 0 {
		levels = []model.Level{defaultLevel()}
	}

	lvl := levels[0]
	cube := buildCube(lvl)

	return &Game{
		cube:       cube,
		renderer:   renderer.New(style),
		sim:        simulation.New(lvl, cube),
		levels:     levels,
		levelIndex: 0,
	}
}

func (g *Game) Update() {
	dt := float64(rl.GetFrameTime())
	sim := g.sim

	g.renderer.Update(g.cube, g.drawState())

	switch sim.State() {
	case simulation.StatePlacement:
		if rl.IsKeyPressed(rl.KeySpace) {
			sim.Start()
		}
		if rl.IsKeyPressed(rl.KeyR) {
			sim.Reset()
		}
		hovered := g.renderer.Hovered()
		if hovered.Valid {
			pos := model.Position{Face: hovered.FaceIdx, Row: hovered.Row, Col: hovered.Col}
			if rl.IsKeyPressed(rl.KeyOne) {
				sim.PlaceMarker(pos, model.MarkerTurnLeft)
			}
			if rl.IsKeyPressed(rl.KeyTwo) {
				sim.PlaceMarker(pos, model.MarkerTurnRight)
			}
			if rl.IsKeyPressed(rl.KeyThree) {
				sim.PlaceMarker(pos, model.MarkerTurnAround)
			}
		}

	case simulation.StateRunning:
		if rl.IsKeyPressed(rl.KeySpace) {
			sim.Stop()
		}
		sim.Update(dt)

	case simulation.StateWon:
		if rl.IsKeyPressed(rl.KeySpace) || rl.IsKeyPressed(rl.KeyEnter) {
			g.nextLevel()
		}
	}
}

func (g *Game) Draw() {
	g.renderer.Draw(g.cube, g.drawState())
}

func (g *Game) Close() {
	g.renderer.Unload()
}

func (g *Game) drawState() renderer.State {
	sim := g.sim
	var phase model.GamePhase
	switch sim.State() {
	case simulation.StatePlacement:
		phase = model.PhasePlacement
	case simulation.StateRunning:
		phase = model.PhaseRunning
	case simulation.StateWon:
		phase = model.PhaseWon
	}
	return renderer.State{
		State: model.State{
			Phase:     phase,
			Pawn:      sim.Pawn(),
			Markers:   sim.Markers(),
			Start:     g.levels[g.levelIndex].Start,
			Goal:      g.levels[g.levelIndex].Goal,
			Inventory: sim.Inventory(),
		},
		PawnPrev:  sim.PrevPawn(),
		PawnAnimT: sim.AnimProgress(),
	}
}

func (g *Game) nextLevel() {
	g.levelIndex++
	if g.levelIndex >= len(g.levels) {
		g.levelIndex = 0
	}
	lvl := g.levels[g.levelIndex]
	rebuildCubeFromLevel(g.cube, lvl)
	g.sim.NextLevel(lvl)
}

func buildCube(lvl model.Level) *model.Cube {
	cube := model.NewCube(lvl.GridSize)
	for _, lc := range lvl.Cells {
		cube.Faces[lc.Position.Face][lc.Position.Row][lc.Position.Col].Kind = lc.Kind
	}
	return cube
}

func rebuildCubeFromLevel(cube *model.Cube, lvl model.Level) {
	n := lvl.GridSize
	cube.GridSize = n
	for i := range cube.Faces {
		cube.Faces[i] = make([][]model.Cell, n)
		for j := range cube.Faces[i] {
			cube.Faces[i][j] = make([]model.Cell, n)
		}
	}
	for _, lc := range lvl.Cells {
		cube.Faces[lc.Position.Face][lc.Position.Row][lc.Position.Col].Kind = lc.Kind
	}
}

func defaultLevel() model.Level {
	return model.Level{
		GridSize:         4,
		Start:            model.Position{Face: 4, Row: 0, Col: 0},
		Goal:             model.Position{Face: 4, Row: 3, Col: 3},
		InitialDirection: model.DirRight,
		Inventory:        model.Inventory{TurnLeft: 1, TurnRight: 1, TurnAround: 1},
	}
}
