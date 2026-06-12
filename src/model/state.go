package model

type GamePhase int

const (
	PhasePlacement GamePhase = iota
	PhaseRunning
	PhaseWon
)

type State struct {
	Phase     GamePhase
	Pawn      Pawn
	Markers   map[Position]MarkerKind
	Start     Position
	Goal      Position
	Inventory Inventory
}
