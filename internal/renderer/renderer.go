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

type Renderer struct {
	camera   rl.Camera3D
	radius   float32
	azimuth  float32
	altitude float32
	style    Style
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

func (r *Renderer) Update() {
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
}

func (r *Renderer) Draw(cube *model.Cube) {
	rl.ClearBackground(r.style.BackgroundColor)
	rl.BeginMode3D(r.camera)
	drawCubeGrid(cube, r.camera.Position, r.style)
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

func drawCubeGrid(cube *model.Cube, cameraPos rl.Vector3, style Style) {
	s := float32(1.0)
	n := cube.GridSize
	step := 2 * s / float32(n)

	faces := [6]faceDesc{
		{rl.NewVector3(s, -s, s), rl.NewVector3(0, 0, -step), rl.NewVector3(0, step, 0), rl.NewVector3(1, 0, 0)},   // +X
		{rl.NewVector3(-s, -s, -s), rl.NewVector3(0, 0, step), rl.NewVector3(0, step, 0), rl.NewVector3(-1, 0, 0)}, // -X
		{rl.NewVector3(-s, s, -s), rl.NewVector3(0, 0, step), rl.NewVector3(step, 0, 0), rl.NewVector3(0, 1, 0)},   // +Y  du/dv swapped: du×dv = (0,1,0)
		{rl.NewVector3(-s, -s, s), rl.NewVector3(0, 0, -step), rl.NewVector3(step, 0, 0), rl.NewVector3(0, -1, 0)}, // -Y  du/dv swapped: du×dv = (0,-1,0)
		{rl.NewVector3(-s, -s, s), rl.NewVector3(step, 0, 0), rl.NewVector3(0, step, 0), rl.NewVector3(0, 0, 1)},   // +Z
		{rl.NewVector3(s, -s, -s), rl.NewVector3(-step, 0, 0), rl.NewVector3(0, step, 0), rl.NewVector3(0, 0, -1)}, // -Z
	}

	var front []faceDesc
	for _, f := range faces {
		if f.normal.X*cameraPos.X+f.normal.Y*cameraPos.Y+f.normal.Z*cameraPos.Z > 0 {
			front = append(front, f)
		}
	}
	for _, f := range front {
		drawFace(f, n, style)
	}
	for _, f := range front {
		drawFaceGrid(f, n, style)
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
