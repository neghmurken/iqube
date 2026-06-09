package model

const (
	FacePosX = 0
	FaceNegX = 1
	FacePosY = 2
	FaceNegY = 3
	FacePosZ = 4
	FaceNegZ = 5
)

type Cell struct{}

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
