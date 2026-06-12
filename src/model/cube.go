package model

const (
	RG = 0
	LF = 1
	UP = 2
	BT = 3
	FR = 4
	BK = 5
)

type CellKind int

const (
	CellNormal  CellKind = iota
	CellVoid             // pawn cannot pass; simulation stops, back to placement
	CellBlocked          // pawn cannot pass; stays on previous cell, simulation continues
)

type Cell struct {
	Kind CellKind
}

type Cube struct {
	GridSize int
	Faces    [6][][]Cell
}

func NewCube(gridSize int) *Cube {
	c := &Cube{GridSize: gridSize}
	for i := range c.Faces {
		c.Faces[i] = make([][]Cell, gridSize)
		for j := range c.Faces[i] {
			c.Faces[i][j] = make([]Cell, gridSize)
		}
	}
	return c
}
