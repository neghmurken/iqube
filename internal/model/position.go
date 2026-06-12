package model

type Position struct {
	Face int
	Row  int
	Col  int
}

func (p Position) Equal(other Position) bool {
	return p.Face == other.Face && p.Row == other.Row && p.Col == other.Col
}
