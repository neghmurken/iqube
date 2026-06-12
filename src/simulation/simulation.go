package simulation

import (
	"github.com/neghmurken/iqube/src/model"
)

type State int

const (
	StatePlacement State = iota
	StateRunning
	StateWon
)

const tickInterval = 0.3

type Simulation struct {
	level       model.Level
	cube        *model.Cube
	pawn        model.Pawn
	prevPawn    model.Pawn
	markers     map[model.Position]model.MarkerKind
	inventory   model.Inventory
	state       State
	accumulator float64
}

func New(level model.Level, cube *model.Cube) *Simulation {
	s := &Simulation{
		level:   level,
		cube:    cube,
		markers: make(map[model.Position]model.MarkerKind),
	}
	s.resetPawn()
	s.inventory = level.Inventory
	return s
}

func (s *Simulation) State() State                                 { return s.state }
func (s *Simulation) Pawn() model.Pawn                             { return s.pawn }
func (s *Simulation) PrevPawn() model.Pawn                         { return s.prevPawn }
func (s *Simulation) AnimProgress() float32                        { return float32(s.accumulator / tickInterval) }
func (s *Simulation) Markers() map[model.Position]model.MarkerKind { return s.markers }
func (s *Simulation) Inventory() model.Inventory                   { return s.inventory }

func (s *Simulation) Start() {
	if s.state == StatePlacement {
		s.state = StateRunning
		s.accumulator = 0
	}
}

func (s *Simulation) Stop() {
	if s.state == StateRunning {
		s.resetPawn()
		s.state = StatePlacement
	}
}

func (s *Simulation) Reset() {
	if s.state == StatePlacement {
		s.markers = make(map[model.Position]model.MarkerKind)
		s.inventory = s.level.Inventory
		s.resetPawn()
	}
}

func (s *Simulation) NextLevel(level model.Level) {
	s.level = level
	s.markers = make(map[model.Position]model.MarkerKind)
	s.inventory = level.Inventory
	s.state = StatePlacement
	s.rebuildCube()
	s.resetPawn()
}

func (s *Simulation) PlaceMarker(pos model.Position, kind model.MarkerKind) {
	if s.state != StatePlacement {
		return
	}
	if existing, ok := s.markers[pos]; ok {
		if existing == kind {
			s.removeMarkerAt(pos)
			return
		}
		s.returnToInventory(existing)
	}
	if !s.spendFromInventory(kind) {
		return
	}
	s.markers[pos] = kind
}

func (s *Simulation) Update(dt float64) {
	if s.state != StateRunning {
		return
	}
	s.accumulator += dt
	for s.accumulator >= tickInterval {
		s.accumulator -= tickInterval
		s.step()
		if s.state != StateRunning {
			return
		}
	}
}

func (s *Simulation) step() {
	s.prevPawn = s.pawn

	if s.cube.Faces[s.pawn.Position.Face][s.pawn.Position.Row][s.pawn.Position.Col].Kind == model.CellVoid {
		s.resetPawn()
		s.state = StatePlacement
		return
	}

	next := s.advance(s.pawn)

	cell := s.cube.Faces[next.Position.Face][next.Position.Row][next.Position.Col]
	switch cell.Kind {
	case model.CellVoid:
		s.pawn = next
		return
	case model.CellBlocked:
		return
	}

	s.pawn = next

	if s.pawn.Position.Equal(s.level.Goal) {
		s.state = StateWon
		return
	}

	if marker, ok := s.markers[s.pawn.Position]; ok {
		switch marker {
		case model.MarkerTurnLeft:
			s.pawn.Direction = s.pawn.Direction.TurnLeft()
		case model.MarkerTurnRight:
			s.pawn.Direction = s.pawn.Direction.TurnRight()
		case model.MarkerTurnAround:
			s.pawn.Direction = s.pawn.Direction.Reverse()
		}
	}
}

func (s *Simulation) advance(p model.Pawn) model.Pawn {
	n := s.level.GridSize
	row, col, face, dir := p.Position.Row, p.Position.Col, p.Position.Face, p.Direction

	switch dir {
	case model.DirUp:
		row++
	case model.DirDown:
		row--
	case model.DirLeft:
		col--
	case model.DirRight:
		col++
	}

	if row >= 0 && row < n && col >= 0 && col < n {
		return model.Pawn{Position: model.Position{Face: face, Row: row, Col: col}, Direction: dir}
	}

	t := model.Transitions[face][dir]
	newRow, newCol := t.Transform(p.Position.Row, p.Position.Col, n)
	return model.Pawn{
		Position:  model.Position{Face: t.Face, Row: newRow, Col: newCol},
		Direction: t.Direction,
	}
}

func (s *Simulation) resetPawn() {
	s.pawn = model.Pawn{
		Position:  s.level.Start,
		Direction: s.level.InitialDirection,
	}
	s.prevPawn = s.pawn
}

func (s *Simulation) rebuildCube() {
	n := s.level.GridSize
	s.cube.GridSize = n
	for i := range s.cube.Faces {
		s.cube.Faces[i] = make([][]model.Cell, n)
		for j := range s.cube.Faces[i] {
			s.cube.Faces[i][j] = make([]model.Cell, n)
		}
	}
	for _, lc := range s.level.Cells {
		s.cube.Faces[lc.Position.Face][lc.Position.Row][lc.Position.Col].Kind = lc.Kind
	}
}

func (s *Simulation) removeMarkerAt(pos model.Position) {
	if existing, ok := s.markers[pos]; ok {
		s.returnToInventory(existing)
		delete(s.markers, pos)
	}
}

func (s *Simulation) returnToInventory(kind model.MarkerKind) {
	switch kind {
	case model.MarkerTurnLeft:
		s.inventory.TurnLeft++
	case model.MarkerTurnRight:
		s.inventory.TurnRight++
	case model.MarkerTurnAround:
		s.inventory.TurnAround++
	}
}

func (s *Simulation) spendFromInventory(kind model.MarkerKind) bool {
	switch kind {
	case model.MarkerTurnLeft:
		if s.inventory.TurnLeft <= 0 {
			return false
		}
		s.inventory.TurnLeft--
	case model.MarkerTurnRight:
		if s.inventory.TurnRight <= 0 {
			return false
		}
		s.inventory.TurnRight--
	case model.MarkerTurnAround:
		if s.inventory.TurnAround <= 0 {
			return false
		}
		s.inventory.TurnAround--
	}
	return true
}
