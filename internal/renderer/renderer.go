package renderer

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/neghmurken/iqube/internal/model"
)

type HoveredCell struct {
	FaceIdx int
	Row     int
	Col     int
	Valid   bool
}

type State struct {
	model.State
	PawnPrev  model.Pawn
	PawnAnimT float32
}

type Renderer struct {
	camera        rl.Camera3D
	radius        float32
	azimuth       float32
	altitude      float32
	style         Style
	faces         [6]faceDesc
	facesGridSize int
	hovered       HoveredCell
	font          rl.Font
	markerTexture rl.Texture2D
}

func New(style Style) *Renderer {
	r := &Renderer{
		radius:   float32(math.Sqrt(20)),
		azimuth:  float32(math.Pi / 4),
		altitude: float32(math.Asin(1.0 / math.Sqrt(3))),
		style:    style,
	}
	r.camera = rl.Camera3D{
		Target:     rl.NewVector3(0, 0, 0),
		Up:         rl.NewVector3(0, 1, 0),
		Fovy:       70,
		Projection: rl.CameraPerspective,
	}
	r.syncCameraPosition()
	if style.FontPath != "" {
		r.font = rl.LoadFontEx(style.FontPath, style.FontSize, nil, 0)
	} else {
		r.font = rl.GetFontDefault()
	}
	if style.MarkerSpritePath != "" {
		r.markerTexture = rl.LoadTexture(style.MarkerSpritePath)
	}
	return r
}

func (r *Renderer) Unload() {
	if r.style.FontPath != "" {
		rl.UnloadFont(r.font)
	}
	if r.markerTexture.ID != 0 {
		rl.UnloadTexture(r.markerTexture)
	}
}

func (r *Renderer) Hovered() HoveredCell { return r.hovered }

func (r *Renderer) Update(cube *model.Cube, ds State) {
	if cube.GridSize != r.facesGridSize {
		r.faces = buildFaces(cube.GridSize)
		r.facesGridSize = cube.GridSize
	}
	r.handleCameraInput()
	h := r.hitTest(cube)
	if h.Valid {
		pos := model.Position{Face: h.FaceIdx, Row: h.Row, Col: h.Col}
		if ds.Start.Equal(pos) || ds.Goal.Equal(pos) || cube.Faces[pos.Face][pos.Row][pos.Col].Kind == model.CellFilled {
			h = HoveredCell{}
		}
	}
	r.hovered = h
}

func (r *Renderer) Draw(cube *model.Cube, ds State) {
	rl.ClearBackground(r.style.GetColor("background"))
	rl.BeginMode3D(r.camera)
	r.drawCubeGrid(cube, r.camera.Position, ds)
	rl.EndMode3D()
	r.drawHUD(ds)
}
