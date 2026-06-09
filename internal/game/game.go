package game

import (
	"github.com/neghmurken/iqube/internal/model"
	"github.com/neghmurken/iqube/internal/renderer"
)

type Game struct {
	cube     *model.Cube
	renderer *renderer.Renderer
}

func New() *Game {
	style, err := renderer.LoadStyle("assets/theme.json")
	if err != nil {
		style = renderer.DefaultStyle()
	}
	return &Game{
		cube:     model.NewCube(4),
		renderer: renderer.New(style),
	}
}

func (g *Game) Update() {
	g.renderer.Update(g.cube)
}

func (g *Game) Draw() {
	g.renderer.Draw(g.cube)
}
