package renderer

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"github.com/neghmurken/iqube/src/model"
)

func (r *Renderer) cellColor(pos model.Position, cube *model.Cube, ds State) (rl.Color, bool) {
	if ds.Goal.Equal(pos) {
		return r.style.GetColor("goal"), true
	}
	if ds.Start.Equal(pos) {
		return r.style.GetColor("start"), true
	}
	switch cube.Faces[pos.Face][pos.Row][pos.Col].Kind {
	case model.CellVoid:
		return r.style.GetColor("void"), true
	case model.CellBlocked:
		return r.style.GetColor("filled"), true
	}
	return rl.Color{}, false
}

func (r *Renderer) cellWorldCenter(pos model.Position) rl.Vector3 {
	f := r.faces[pos.Face]
	fc := float32(pos.Col) + 0.5
	fr := float32(pos.Row) + 0.5
	return rl.NewVector3(
		f.p0.X+f.du.X*fc+f.dv.X*fr,
		f.p0.Y+f.du.Y*fc+f.dv.Y*fr,
		f.p0.Z+f.du.Z*fc+f.dv.Z*fr,
	)
}

func (r *Renderer) drawPawn(ds State, color rl.Color) {
	from3D := r.cellWorldCenter(ds.PawnPrev.Position)
	to3D := r.cellWorldCenter(ds.Pawn.Position)
	t := ds.PawnAnimT

	var pos, normal rl.Vector3
	fn0 := r.faces[ds.PawnPrev.Position.Face].normal
	fn1 := r.faces[ds.Pawn.Position.Face].normal

	if ds.PawnPrev.Position.Face == ds.Pawn.Position.Face {
		pos = rl.Vector3Lerp(from3D, to3D, t)
		normal = fn1
	} else {
		pick := func(n0, n1, p float32) float32 {
			if n0 != 0 {
				return n0
			}
			if n1 != 0 {
				return n1
			}
			return p
		}
		edge := rl.NewVector3(
			pick(fn0.X, fn1.X, from3D.X),
			pick(fn0.Y, fn1.Y, from3D.Y),
			pick(fn0.Z, fn1.Z, from3D.Z),
		)
		if t < 0.5 {
			pos = rl.Vector3Lerp(from3D, edge, t*2)
		} else {
			pos = rl.Vector3Lerp(edge, to3D, (t-0.5)*2)
		}
		ln := rl.Vector3Lerp(fn0, fn1, t)
		if rl.Vector3Length(ln) > 0 {
			normal = rl.Vector3Normalize(ln)
		} else {
			normal = fn1
		}
	}

	f := r.faces[ds.Pawn.Position.Face]
	radius := rl.Vector3Length(f.du) * 0.2
	off := radius + lineOffset
	p := rl.NewVector3(pos.X+normal.X*off, pos.Y+normal.Y*off, pos.Z+normal.Z*off)
	rl.DrawSphere(p, radius, color)
}

func (r *Renderer) drawMarker(f faceDesc, row, col int, kind model.MarkerKind) {
	if r.markerTexture.ID == 0 {
		return
	}

	texW := float32(r.markerTexture.Width)
	texH := float32(r.markerTexture.Height)
	hw, hh := texW/2, texH/2

	var srcX, srcY float32
	switch kind {
	case model.MarkerTurnAround:
		srcX, srcY = 0, 0
	case model.MarkerTurnRight:
		srcX, srcY = 0, hh
	case model.MarkerTurnLeft:
		srcX, srcY = hw, hh
	}

	u0 := srcX / texW
	u1 := (srcX + hw) / texW
	v0 := srcY / texH
	v1 := (srcY + hh) / texH

	fc := float32(col)
	fr := float32(row)
	ox := f.normal.X * lineOffset * 5
	oy := f.normal.Y * lineOffset * 5
	oz := f.normal.Z * lineOffset * 5

	pt := func(u, v float32) rl.Vector3 {
		return rl.NewVector3(
			f.p0.X+f.du.X*u+f.dv.X*v+ox,
			f.p0.Y+f.du.Y*u+f.dv.Y*v+oy,
			f.p0.Z+f.du.Z*u+f.dv.Z*v+oz,
		)
	}

	p0 := pt(fc, fr)
	p1 := pt(fc+1, fr)
	p2 := pt(fc+1, fr+1)
	p3 := pt(fc, fr+1)

	rl.SetTexture(r.markerTexture.ID)
	rl.Begin(rl.Quads)
	rl.Color4ub(255, 255, 255, 255)
	rl.Normal3f(f.normal.X, f.normal.Y, f.normal.Z)
	rl.TexCoord2f(u0, v1)
	rl.Vertex3f(p0.X, p0.Y, p0.Z)
	rl.TexCoord2f(u1, v1)
	rl.Vertex3f(p1.X, p1.Y, p1.Z)
	rl.TexCoord2f(u1, v0)
	rl.Vertex3f(p2.X, p2.Y, p2.Z)
	rl.TexCoord2f(u0, v0)
	rl.Vertex3f(p3.X, p3.Y, p3.Z)
	rl.End()
	rl.SetTexture(0)
}
