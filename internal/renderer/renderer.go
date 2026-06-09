package renderer

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/neghmurken/iqube/internal/model"
)

const (
	dragSensitivity = 0.005
	zoomSensitivity = 0.5
	minRadius       = 2.5
	maxRadius       = 30.0
)

type hoveredCell struct {
	faceIdx int
	row     int
	col     int
	valid   bool
}

type Renderer struct {
	camera        rl.Camera3D
	radius        float32
	azimuth       float32
	altitude      float32
	style         Style
	faces         [6]faceDesc
	facesGridSize int
	hovered       hoveredCell
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
	return r
}

func (r *Renderer) Update(cube *model.Cube) {
	if cube.GridSize != r.facesGridSize {
		r.faces = buildFaces(cube.GridSize)
		r.facesGridSize = cube.GridSize
	}
	if wheel := rl.GetMouseWheelMove(); wheel != 0 {
		r.radius = clamp(r.radius-wheel*zoomSensitivity, minRadius, maxRadius)
		r.syncCameraPosition()
	}
	if rl.IsMouseButtonDown(rl.MouseButtonRight) {
		delta := rl.GetMouseDelta()
		r.azimuth += delta.X * dragSensitivity
		r.altitude += delta.Y * dragSensitivity
		r.altitude = clamp(r.altitude, float32(-math.Pi/2)+0.01, float32(math.Pi/2)-0.01)
		r.syncCameraPosition()
	}
	r.hovered = r.hitTest(cube)
}

func (r *Renderer) hitTest(cube *model.Cube) hoveredCell {
	ray := rl.GetScreenToWorldRay(rl.GetMousePosition(), r.camera)
	n := cube.GridSize
	fn := float32(n)

	best := hoveredCell{}
	bestDist := float32(math.MaxFloat32)

	for faceIdx, f := range r.faces {
		if f.normal.X*r.camera.Position.X+f.normal.Y*r.camera.Position.Y+f.normal.Z*r.camera.Position.Z <= 0 {
			continue
		}
		v0 := f.p0
		v1 := rl.NewVector3(f.p0.X+f.du.X*fn, f.p0.Y+f.du.Y*fn, f.p0.Z+f.du.Z*fn)
		v2 := rl.NewVector3(f.p0.X+f.du.X*fn+f.dv.X*fn, f.p0.Y+f.du.Y*fn+f.dv.Y*fn, f.p0.Z+f.du.Z*fn+f.dv.Z*fn)
		v3 := rl.NewVector3(f.p0.X+f.dv.X*fn, f.p0.Y+f.dv.Y*fn, f.p0.Z+f.dv.Z*fn)

		c1 := rl.GetRayCollisionTriangle(ray, v0, v1, v2)
		c2 := rl.GetRayCollisionTriangle(ray, v0, v2, v3)
		var col rl.RayCollision
		if c1.Hit && (!c2.Hit || c1.Distance < c2.Distance) {
			col = c1
		} else if c2.Hit {
			col = c2
		}
		if !col.Hit || col.Distance >= bestDist {
			continue
		}

		lx := col.Point.X - f.p0.X
		ly := col.Point.Y - f.p0.Y
		lz := col.Point.Z - f.p0.Z
		duLen2 := f.du.X*f.du.X + f.du.Y*f.du.Y + f.du.Z*f.du.Z
		dvLen2 := f.dv.X*f.dv.X + f.dv.Y*f.dv.Y + f.dv.Z*f.dv.Z
		u := (lx*f.du.X + ly*f.du.Y + lz*f.du.Z) / duLen2
		v := (lx*f.dv.X + ly*f.dv.Y + lz*f.dv.Z) / dvLen2

		cellCol := max(0, min(n-1, int(u)))
		cellRow := max(0, min(n-1, int(v)))
		bestDist = col.Distance
		best = hoveredCell{faceIdx: faceIdx, row: cellRow, col: cellCol, valid: true}
	}
	return best
}

func (r *Renderer) Draw(cube *model.Cube) {
	rl.ClearBackground(r.style.BackgroundColor)
	rl.BeginMode3D(r.camera)
	r.drawCubeGrid(cube, r.camera.Position, r.style)
	rl.EndMode3D()
}

func (r *Renderer) syncCameraPosition() {
	cosAlt := float32(math.Cos(float64(r.altitude)))
	sinAlt := float32(math.Sin(float64(r.altitude)))
	cosAz := float32(math.Cos(float64(r.azimuth)))
	sinAz := float32(math.Sin(float64(r.azimuth)))
	r.camera.Position = rl.NewVector3(
		r.radius*cosAlt*cosAz,
		r.radius*sinAlt,
		r.radius*cosAlt*sinAz,
	)
}

func clamp(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

type faceDesc struct {
	p0     rl.Vector3
	du     rl.Vector3
	dv     rl.Vector3
	normal rl.Vector3
}

const lineOffset = float32(0.002)

func buildFaces(gridSize int) [6]faceDesc {
	s := float32(1.0)
	step := 2 * s / float32(gridSize)
	return [6]faceDesc{
		{rl.NewVector3(s, -s, s), rl.NewVector3(0, 0, -step), rl.NewVector3(0, step, 0), rl.NewVector3(1, 0, 0)},   // +X
		{rl.NewVector3(-s, -s, -s), rl.NewVector3(0, 0, step), rl.NewVector3(0, step, 0), rl.NewVector3(-1, 0, 0)}, // -X
		{rl.NewVector3(-s, s, -s), rl.NewVector3(0, 0, step), rl.NewVector3(step, 0, 0), rl.NewVector3(0, 1, 0)},   // +Y
		{rl.NewVector3(-s, -s, s), rl.NewVector3(0, 0, -step), rl.NewVector3(step, 0, 0), rl.NewVector3(0, -1, 0)}, // -Y
		{rl.NewVector3(-s, -s, s), rl.NewVector3(step, 0, 0), rl.NewVector3(0, step, 0), rl.NewVector3(0, 0, 1)},   // +Z
		{rl.NewVector3(s, -s, -s), rl.NewVector3(-step, 0, 0), rl.NewVector3(0, step, 0), rl.NewVector3(0, 0, -1)}, // -Z
	}
}

func (r *Renderer) drawCubeGrid(cube *model.Cube, cameraPos rl.Vector3, style Style) {
	var front [6]faceDesc
	count := 0
	for _, f := range r.faces {
		if f.normal.X*cameraPos.X+f.normal.Y*cameraPos.Y+f.normal.Z*cameraPos.Z > 0 {
			front[count] = f
			count++
		}
	}
	for i := range count {
		drawFace(front[i], cube.GridSize, style)
	}
	if r.hovered.valid {
		drawHoveredCell(r.faces[r.hovered.faceIdx], r.hovered.row, r.hovered.col, style.HoverColor)
	}
	for i := range count {
		drawFaceGrid(front[i], cube.GridSize, style)
	}
}

func drawFace(f faceDesc, n int, style Style) {
	fn := float32(n)
	v0 := f.p0
	v1 := rl.NewVector3(f.p0.X+f.du.X*fn, f.p0.Y+f.du.Y*fn, f.p0.Z+f.du.Z*fn)
	v2 := rl.NewVector3(f.p0.X+f.du.X*fn+f.dv.X*fn, f.p0.Y+f.du.Y*fn+f.dv.Y*fn, f.p0.Z+f.du.Z*fn+f.dv.Z*fn)
	v3 := rl.NewVector3(f.p0.X+f.dv.X*fn, f.p0.Y+f.dv.Y*fn, f.p0.Z+f.dv.Z*fn)
	rl.DrawTriangle3D(v0, v1, v2, style.FaceColor)
	rl.DrawTriangle3D(v0, v2, v3, style.FaceColor)
}

func drawHoveredCell(f faceDesc, row, col int, color rl.Color) {
	fr := float32(row)
	fc := float32(col)
	ox := f.normal.X * lineOffset * 2
	oy := f.normal.Y * lineOffset * 2
	oz := f.normal.Z * lineOffset * 2
	v0 := rl.NewVector3(f.p0.X+f.du.X*fc+f.dv.X*fr+ox, f.p0.Y+f.du.Y*fc+f.dv.Y*fr+oy, f.p0.Z+f.du.Z*fc+f.dv.Z*fr+oz)
	v1 := rl.NewVector3(v0.X+f.du.X, v0.Y+f.du.Y, v0.Z+f.du.Z)
	v2 := rl.NewVector3(v0.X+f.du.X+f.dv.X, v0.Y+f.du.Y+f.dv.Y, v0.Z+f.du.Z+f.dv.Z)
	v3 := rl.NewVector3(v0.X+f.dv.X, v0.Y+f.dv.Y, v0.Z+f.dv.Z)
	rl.DrawTriangle3D(v0, v1, v2, color)
	rl.DrawTriangle3D(v0, v2, v3, color)
}

func drawFaceGrid(f faceDesc, n int, style Style) {
	fn := float32(n)
	ox := f.normal.X * lineOffset
	oy := f.normal.Y * lineOffset
	oz := f.normal.Z * lineOffset

	for i := 0; i <= n; i++ {
		fi := float32(i)

		startU := rl.NewVector3(f.p0.X+f.du.X*fi+ox, f.p0.Y+f.du.Y*fi+oy, f.p0.Z+f.du.Z*fi+oz)
		endU := rl.NewVector3(startU.X+f.dv.X*fn, startU.Y+f.dv.Y*fn, startU.Z+f.dv.Z*fn)
		rl.DrawLine3D(startU, endU, style.GridColor)

		startV := rl.NewVector3(f.p0.X+f.dv.X*fi+ox, f.p0.Y+f.dv.Y*fi+oy, f.p0.Z+f.dv.Z*fi+oz)
		endV := rl.NewVector3(startV.X+f.du.X*fn, startV.Y+f.du.Y*fn, startV.Z+f.du.Z*fn)
		rl.DrawLine3D(startV, endV, style.GridColor)
	}
}
