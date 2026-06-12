# 002 — Pawn movement animation

Animate the pawn smoothly between cells instead of jumping cell-to-cell.

## Context

- Language: Go. Renderer: Raylib via `github.com/gen2brain/raylib-go/raylib`.
- Simulation ticks at `tickInterval = 0.3s` (`internal/simulation/simulation.go:15`).
- Each tick, `step()` teleports `s.pawn` to the next cell instantly.
- `accumulator` tracks elapsed time since the last tick; it resets by subtracting `tickInterval` after each step.
- The renderer draws the pawn at the exact discrete cell position every frame.

## Goal

Lerp the pawn's **3D world position** between the previous cell center and the current cell center using `t = accumulator / tickInterval` (range [0, 1)).

Direction (sprite rotation) always reflects `pawn.Direction` — the current logical direction — not the travel direction. Do not interpolate direction.

## Files to change

### `internal/simulation/simulation.go`

1. Add `prevPawn model.Pawn` field to `Simulation`.
2. In `step()`, before advancing: `s.prevPawn = s.pawn`.
3. In `resetPawn()`, after setting `s.pawn`: `s.prevPawn = s.pawn`.
4. Add getter: `func (s *Simulation) PrevPawn() model.Pawn { return s.prevPawn }`.
5. Add getter: `func (s *Simulation) AnimProgress() float32 { return float32(s.accumulator / tickInterval) }`.

### `internal/renderer/renderer.go`

`DrawState` (defined here):

```go
type DrawState struct {
    Phase     GamePhase
    Pawn      model.Pawn
    Markers   map[model.Position]model.MarkerKind
    Start     model.Position
    Goal      model.Position
    Inventory model.Inventory
}
```

Add two fields:

```go
PawnPrev  model.Pawn
PawnAnimT float32
```

Add a helper method on `Renderer`:

```go
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
```

Rewrite `drawPawn` to:
- Compute `from3D = r.cellWorldCenter(ds.PawnPrev.Position)` and `to3D = r.cellWorldCenter(ds.Pawn.Position)`.
- Lerp: `pos = lerp(from3D, to3D, ds.PawnAnimT)`.
- Draw the sprite quad centered on `pos + normal*offset`, aligned to the face of `ds.Pawn.Position` (current face — use `r.faces[ds.Pawn.Position.Face]`).
- Rotation from `ds.Pawn.Direction` (unchanged logic from current `drawPawn`).

Current `drawPawn` signature: `func (r *Renderer) drawPawn(f faceDesc, row, col int, dir model.Direction)`.

Change the call site in `drawCubeGrid` (currently iterates over `frontIdx` to find the pawn's face) to pass `ds` instead, and call the new signature directly with `ds`.

Current call site (renderer.go ~line 261):
```go
for i := range count {
    if frontIdx[i] == ds.Pawn.Position.Face {
        r.drawPawn(front[i], ds.Pawn.Position.Row, ds.Pawn.Position.Col, ds.Pawn.Direction)
    }
}
```

The face guard is still needed (don't draw the pawn on a back-facing face), but the quad center is now the lerped 3D position rather than the cell corners.

### `internal/game/game.go`

In `drawState()`, populate the two new fields:

```go
PawnPrev:  sim.PrevPawn(),
PawnAnimT: sim.AnimProgress(),
```

`AnimProgress()` returns 0 when not running (accumulator is 0), so no special casing needed for `StatePlacement` / `StateWon`.

## Constraints

- Animation only affects the rendered position; simulation logic is unchanged.
- On `Stop()` / `resetPawn()` the pawn snaps back to start — `prevPawn = pawn` ensures `t` has no effect (from == to → lerp is a no-op).
- `r.faces` is valid during `Draw` because `Update` (which calls `buildFaces` when grid size changes) always runs before `Draw`.
- Do not change `tickInterval` or simulation timing.
- Build with `make build`.
