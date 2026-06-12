package renderer

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/neghmurken/iqube/internal/model"
)

func (r *Renderer) drawCubeGrid(cube *model.Cube, cameraPos rl.Vector3, ds State) {
	var front [6]faceDesc
	var frontIdx [6]int
	count := 0
	for i, f := range r.faces {
		if f.normal.X*cameraPos.X+f.normal.Y*cameraPos.Y+f.normal.Z*cameraPos.Z > 0 {
			front[count] = f
			frontIdx[count] = i
			count++
		}
	}

	for i := range count {
		drawFace(front[i], cube.GridSize, r.style.GetColor("face"))
	}

	for i := range count {
		fi := frontIdx[i]
		n := cube.GridSize
		for row := range n {
			for col := range n {
				pos := model.Position{Face: fi, Row: row, Col: col}
				color, ok := r.cellColor(pos, cube, ds)
				if ok {
					drawCell(front[i], row, col, color)
				}
			}
		}
	}

	if ds.Phase == model.PhasePlacement && r.hovered.Valid {
		drawCell(r.faces[r.hovered.FaceIdx], r.hovered.Row, r.hovered.Col, r.style.GetColor("hover"))
	}

	for pos, kind := range ds.Markers {
		for i := range count {
			if frontIdx[i] == pos.Face {
				r.drawMarker(front[i], pos.Row, pos.Col, kind)
			}
		}
	}

	for i := range count {
		if frontIdx[i] == ds.Pawn.Position.Face || frontIdx[i] == ds.PawnPrev.Position.Face {
			r.drawPawn(ds)
			break
		}
	}

	for i := range count {
		drawFaceGrid(front[i], cube.GridSize, r.style.GetColor("grid"))
	}

	for i := range count {
		drawFaceOrigin(front[i])
	}
}

func drawFace(f faceDesc, n int, color rl.Color) {
	fn := float32(n)
	v0 := f.p0
	v1 := rl.NewVector3(f.p0.X+f.du.X*fn, f.p0.Y+f.du.Y*fn, f.p0.Z+f.du.Z*fn)
	v2 := rl.NewVector3(f.p0.X+f.du.X*fn+f.dv.X*fn, f.p0.Y+f.du.Y*fn+f.dv.Y*fn, f.p0.Z+f.du.Z*fn+f.dv.Z*fn)
	v3 := rl.NewVector3(f.p0.X+f.dv.X*fn, f.p0.Y+f.dv.Y*fn, f.p0.Z+f.dv.Z*fn)
	rl.DrawTriangle3D(v0, v1, v2, color)
	rl.DrawTriangle3D(v0, v2, v3, color)
}

func drawCell(f faceDesc, row, col int, color rl.Color) {
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

func drawFaceGrid(f faceDesc, n int, color rl.Color) {
	fn := float32(n)
	ox := f.normal.X * lineOffset
	oy := f.normal.Y * lineOffset
	oz := f.normal.Z * lineOffset

	for i := 0; i <= n; i++ {
		fi := float32(i)
		startU := rl.NewVector3(f.p0.X+f.du.X*fi+ox, f.p0.Y+f.du.Y*fi+oy, f.p0.Z+f.du.Z*fi+oz)
		endU := rl.NewVector3(startU.X+f.dv.X*fn, startU.Y+f.dv.Y*fn, startU.Z+f.dv.Z*fn)
		rl.DrawLine3D(startU, endU, color)

		startV := rl.NewVector3(f.p0.X+f.dv.X*fi+ox, f.p0.Y+f.dv.Y*fi+oy, f.p0.Z+f.dv.Z*fi+oz)
		endV := rl.NewVector3(startV.X+f.du.X*fn, startV.Y+f.du.Y*fn, startV.Z+f.du.Z*fn)
		rl.DrawLine3D(startV, endV, color)
	}
}

var faceNames = [6]string{"RG", "LF", "UP", "BT", "FR", "BK"}

func drawFaceOrigin(f faceDesc) {
	size := rl.Vector3Length(f.du) * 0.15
	ox := f.normal.X * lineOffset * 5
	oy := f.normal.Y * lineOffset * 5
	oz := f.normal.Z * lineOffset * 5
	v0 := rl.NewVector3(f.p0.X+ox, f.p0.Y+oy, f.p0.Z+oz)
	v1 := rl.NewVector3(v0.X+f.du.X*size, v0.Y+f.du.Y*size, v0.Z+f.du.Z*size)
	v2 := rl.NewVector3(v0.X+f.du.X*size+f.dv.X*size, v0.Y+f.du.Y*size+f.dv.Y*size, v0.Z+f.du.Z*size+f.dv.Z*size)
	v3 := rl.NewVector3(v0.X+f.dv.X*size, v0.Y+f.dv.Y*size, v0.Z+f.dv.Z*size)
	rl.DrawTriangle3D(v0, v1, v2, rl.Blue)
	rl.DrawTriangle3D(v0, v2, v3, rl.Blue)
}
