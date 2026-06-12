package renderer

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/neghmurken/iqube/internal/model"
)

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
		{rl.NewVector3(s, -s, s), rl.NewVector3(0, 0, -step), rl.NewVector3(0, step, 0), rl.NewVector3(1, 0, 0)},
		{rl.NewVector3(-s, -s, -s), rl.NewVector3(0, 0, step), rl.NewVector3(0, step, 0), rl.NewVector3(-1, 0, 0)},
		{rl.NewVector3(-s, s, -s), rl.NewVector3(0, 0, step), rl.NewVector3(step, 0, 0), rl.NewVector3(0, 1, 0)},
		{rl.NewVector3(-s, -s, s), rl.NewVector3(0, 0, -step), rl.NewVector3(step, 0, 0), rl.NewVector3(0, -1, 0)},
		{rl.NewVector3(-s, -s, s), rl.NewVector3(step, 0, 0), rl.NewVector3(0, step, 0), rl.NewVector3(0, 0, 1)},
		{rl.NewVector3(s, -s, -s), rl.NewVector3(-step, 0, 0), rl.NewVector3(0, step, 0), rl.NewVector3(0, 0, -1)},
	}
}

func (r *Renderer) hitTest(cube *model.Cube) HoveredCell {
	ray := rl.GetScreenToWorldRay(rl.GetMousePosition(), r.camera)
	n := cube.GridSize
	fn := float32(n)

	best := HoveredCell{}
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
		best = HoveredCell{FaceIdx: faceIdx, Row: cellRow, Col: cellCol, Valid: true}
	}
	return best
}
