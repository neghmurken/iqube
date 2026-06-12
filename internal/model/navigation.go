package model

// EdgeTransition describes where the pawn lands when it exits a face edge.
type EdgeTransition struct {
	Face      int
	Direction Direction
	// Transform maps (row, col, gridSize) of the exit cell to the entry cell on the new face.
	Transform func(row, col, n int) (int, int)
}

// Transitions[face][dir] gives the transition when the pawn exits face in dir.
var Transitions [6][4]EdgeTransition

func init() {
	// Convention (matches simulation.go):
	//   DirUp   = row++  → exits high-row boundary → exits in +dv world direction
	//   DirDown = row--  → exits low-row boundary  → exits in -dv world direction
	//   DirLeft = col--  → exits low-col boundary  → exits in -du world direction
	//   DirRight= col++  → exits high-col boundary → exits in +du world direction
	//
	// Face axes (step = 2/gridSize):
	//   RG (+X): p0=(+1,-1,+1)  du=(0,0,-s) dv=(0,+s,0)  col→-Z  row→+Y
	//   LF (-X): p0=(-1,-1,-1)  du=(0,0,+s) dv=(0,+s,0)  col→+Z  row→+Y
	//   UP (+Y): p0=(-1,+1,-1)  du=(0,0,+s) dv=(+s,0,0)  col→+Z  row→+X
	//   BT (-Y): p0=(-1,-1,+1)  du=(0,0,-s) dv=(+s,0,0)  col→-Z  row→+X
	//   FR (+Z): p0=(-1,-1,+1)  du=(+s,0,0) dv=(0,+s,0)  col→+X  row→+Y
	//   BK (-Z): p0=(+1,-1,-1)  du=(-s,0,0) dv=(0,+s,0)  col→-X  row→+Y
	//
	// Entry rule: pawn enters new face at the edge row/col adjacent to the source face,
	// going inward (away from that edge).
	//   entry from high-row edge (row=n-1) → DirDown
	//   entry from low-row  edge (row=0)   → DirUp
	//   entry from high-col edge (col=n-1) → DirLeft
	//   entry from low-col  edge (col=0)   → DirRight

	// --- RG (+X)  exits: DirUp→+Y  DirDown→-Y  DirLeft→+Z  DirRight→-Z ---
	Transitions[RG][DirUp] = EdgeTransition{UP, DirDown,
		func(row, col, n int) (int, int) { return n - 1, n - 1 - col }}
	Transitions[RG][DirDown] = EdgeTransition{BT, DirDown,
		func(row, col, n int) (int, int) { return n - 1, col }}
	Transitions[RG][DirLeft] = EdgeTransition{FR, DirLeft,
		func(row, col, n int) (int, int) { return row, n - 1 }}
	Transitions[RG][DirRight] = EdgeTransition{BK, DirRight,
		func(row, col, n int) (int, int) { return row, 0 }}

	// --- LF (-X)  exits: DirUp→+Y  DirDown→-Y  DirLeft→-Z  DirRight→+Z ---
	Transitions[LF][DirUp] = EdgeTransition{UP, DirUp,
		func(row, col, n int) (int, int) { return 0, col }}
	Transitions[LF][DirDown] = EdgeTransition{BT, DirUp,
		func(row, col, n int) (int, int) { return 0, n - 1 - col }}
	Transitions[LF][DirLeft] = EdgeTransition{BK, DirLeft,
		func(row, col, n int) (int, int) { return row, n - 1 }}
	Transitions[LF][DirRight] = EdgeTransition{FR, DirRight,
		func(row, col, n int) (int, int) { return row, 0 }}

	// --- UP (+Y)  exits: DirUp→+X  DirDown→-X  DirLeft→-Z  DirRight→+Z ---
	Transitions[UP][DirUp] = EdgeTransition{RG, DirDown,
		func(row, col, n int) (int, int) { return n - 1, n - 1 - col }}
	Transitions[UP][DirDown] = EdgeTransition{LF, DirDown,
		func(row, col, n int) (int, int) { return n - 1, col }}
	Transitions[UP][DirLeft] = EdgeTransition{BK, DirDown,
		func(row, col, n int) (int, int) { return n - 1, n - 1 - row }}
	Transitions[UP][DirRight] = EdgeTransition{FR, DirDown,
		func(row, col, n int) (int, int) { return n - 1, row }}

	// --- BT (-Y)  exits: DirUp→+X  DirDown→-X  DirLeft→+Z  DirRight→-Z ---
	Transitions[BT][DirUp] = EdgeTransition{RG, DirUp,
		func(row, col, n int) (int, int) { return 0, col }}
	Transitions[BT][DirDown] = EdgeTransition{LF, DirUp,
		func(row, col, n int) (int, int) { return 0, n - 1 - col }}
	Transitions[BT][DirLeft] = EdgeTransition{FR, DirUp,
		func(row, col, n int) (int, int) { return 0, row }}
	Transitions[BT][DirRight] = EdgeTransition{BK, DirUp,
		func(row, col, n int) (int, int) { return 0, n - 1 - row }}

	// --- FR (+Z)  exits: DirUp→+Y  DirDown→-Y  DirLeft→-X  DirRight→+X ---
	Transitions[FR][DirUp] = EdgeTransition{UP, DirLeft,
		func(row, col, n int) (int, int) { return col, n - 1 }}
	Transitions[FR][DirDown] = EdgeTransition{BT, DirRight,
		func(row, col, n int) (int, int) { return col, 0 }}
	Transitions[FR][DirLeft] = EdgeTransition{LF, DirLeft,
		func(row, col, n int) (int, int) { return row, n - 1 }}
	Transitions[FR][DirRight] = EdgeTransition{RG, DirRight,
		func(row, col, n int) (int, int) { return row, 0 }}

	// --- BK (-Z)  exits: DirUp→+Y  DirDown→-Y  DirLeft→+X  DirRight→-X ---
	Transitions[BK][DirUp] = EdgeTransition{UP, DirRight,
		func(row, col, n int) (int, int) { return n - 1 - col, 0 }}
	Transitions[BK][DirDown] = EdgeTransition{BT, DirLeft,
		func(row, col, n int) (int, int) { return n - 1 - col, n - 1 }}
	Transitions[BK][DirLeft] = EdgeTransition{RG, DirLeft,
		func(row, col, n int) (int, int) { return row, n - 1 }}
	Transitions[BK][DirRight] = EdgeTransition{LF, DirRight,
		func(row, col, n int) (int, int) { return row, 0 }}
}
