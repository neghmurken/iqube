package model

type Direction int

const (
	DirUp Direction = iota
	DirDown
	DirLeft
	DirRight
)

type MarkerKind int

const (
	MarkerTurnLeft MarkerKind = iota
	MarkerTurnRight
	MarkerTurnAround
)

func (d Direction) TurnLeft() Direction {
	switch d {
	case DirUp:
		return DirLeft
	case DirLeft:
		return DirDown
	case DirDown:
		return DirRight
	case DirRight:
		return DirUp
	}
	return d
}

func (d Direction) TurnRight() Direction {
	switch d {
	case DirUp:
		return DirRight
	case DirRight:
		return DirDown
	case DirDown:
		return DirLeft
	case DirLeft:
		return DirUp
	}
	return d
}

func (d Direction) Reverse() Direction {
	switch d {
	case DirUp:
		return DirDown
	case DirDown:
		return DirUp
	case DirLeft:
		return DirRight
	case DirRight:
		return DirLeft
	}
	return d
}

func ParseDirection(s string) Direction {
	switch s {
	case "up":
		return DirUp
	case "down":
		return DirDown
	case "left":
		return DirLeft
	case "right":
		return DirRight
	}
	return DirRight
}
