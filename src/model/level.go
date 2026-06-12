package model

type Inventory struct {
	TurnLeft    int
	TurnRight   int
	TurnAround  int
}

type Level struct {
	GridSize         int
	Start            Position
	Goal             Position
	InitialDirection Direction
	Inventory        Inventory
	Cells            []LevelCell
}

type LevelCell struct {
	Position Position
	Kind     CellKind
}
