# iqube — First Implementation

**Date**: 2026-06-10  
**Project**: `/home/pib/projects/iqube`  
**Stack**: Go + Raylib  

---

## Context

iqube is a 3D puzzle game. A pawn moves across faces of a cube (each face = grid of cells). Player places directional markers on cells; pawn follows them. Goal: guide pawn from start cell to goal cell.

Full spec: `/home/pib/projects/iqube/SPECS.md`  
Project instructions: `/home/pib/projects/iqube/CLAUDE.md`

---

## Current State

### What exists (working)
- 3D cube rendering with orbit camera (right-drag, scroll zoom)
- Mouse hover hit-testing on cells
- Theme loading from `assets/theme.json`
- Face coordinate systems defined in `buildFaces()` (`internal/renderer/renderer.go`)

### What does NOT exist yet
- Game logic: pawn, markers, simulation
- Level loading
- Any UI/HUD

### File structure
```
main.go
internal/
  game/game.go          Game struct, owns Cube + Renderer
  model/cube.go         Cube, Cell (empty struct), face constants
  renderer/renderer.go  3D render, camera, hit-test
  renderer/style.go     Style, theme loader
assets/theme.json
```

---

## Architecture Decisions (from design session)

### Data model

| Decision | Choice |
|----------|--------|
| Cell types | `CellKind` enum: `Empty \| Filled` on `Cell` struct |
| Start/Goal | Separate `Position{Face, Row, Col}` fields on `Level` — not on cells |
| Level format | JSON files in `assets/levels/`, loaded at runtime |
| Infinite loop | No detection — player's responsibility |
| Simulation tick | 0.3s per step, time-accumulator based |
| Pawn initial direction | Per-level, defined in JSON |

### Level JSON schema
```json
{
  "grid_size": 4,
  "start": { "face": 4, "row": 0, "col": 0 },
  "goal":  { "face": 2, "row": 3, "col": 3 },
  "initial_direction": "right",
  "inventory": { "turn_left": 2, "turn_right": 1, "turn_around": 1 },
  "cells": [
    { "face": 4, "row": 1, "col": 2, "kind": "filled" }
  ]
}
```
- `cells` sparse — only non-empty entries listed
- `initial_direction`: `up|down|left|right` relative to face local grid (right=+col, down=+row)
- Face indices: 0=+X, 1=-X, 2=+Y, 3=-Y, 4=+Z, 5=-Z (matches existing constants)

### Input scheme
| Key | Action | Constraint |
|-----|--------|-----------|
| `Space` | Start simulation / stop simulation (pawn resets, markers+inventory kept) | — |
| `R` | Reset board: clear markers, restore full inventory, pawn to start | Placement mode only |
| `1` / `2` / `3` | Place marker (TurnLeft/TurnRight/TurnAround) on hovered cell | Placement mode only |
| same hotkey on existing marker | Remove marker | — |
| different hotkey on existing marker | Replace marker (old returned to inventory, new spent) | — |

### State machine
```
Placement → [Space]        → Running
Running   → [Space]        → Placement  (pawn reset, markers/inventory kept)
Running   → [goal reached] → Won
Running   → [hit filled]   → Placement  (auto, markers/inventory kept)
Won       → [any key]      → Placement  (load next level)
Placement → [R]            → Placement  (full reset)
```

### Face edge traversal
Direction is preserved as a 3D world vector when crossing face edges. On new face, project the 3D direction onto the new face's `du`/`dv` axes (already defined in `buildFaces()`). Derive a static `[6][4]EdgeTransition` table from these vectors once at init.

---

## Implementation Plan

Build in this order:

### 1. `internal/model/` — new files
- `direction.go`: `Direction` (Up/Down/Left/Right), `MarkerKind` (TurnLeft/TurnRight/TurnAround)
- `position.go`: `Position{Face, Row, Col}`
- `level.go`: `Level{GridSize, Start, Goal, InitialDirection, Inventory, Cells}`
- `pawn.go`: `Pawn{Position, Direction}`
- `navigation.go`: `EdgeTransition` table derived from `buildFaces()` face descriptors
- Update `cube.go`: add `CellKind` to `Cell`

### 2. `internal/level/loader.go`
- Read `assets/levels/*.json` → `[]Level`
- JSON tags matching schema above

### 3. `internal/simulation/simulation.go`
- State machine (Placement/Running/Won/Lost)
- `Step()` — advance pawn, read marker, check goal/filled
- `PlaceMarker()`, `RemoveMarker()`
- `Start()`, `Stop()`, `Reset()`
- Tick accumulator at 0.3s

### 4. `internal/game/game.go` — update
- Wire simulation + input handling
- Load levels from `level.Loader`

### 5. `internal/renderer/` — additions
- Draw pawn on current cell
- Draw markers on cells (visual indicator per type)
- Highlight start/goal cells
- HUD: inventory counts, current mode, win/loss state

### 6. `assets/levels/001.json`
- First playable level

---

## Suggested Skills

- `/run` — after implementing simulation, use to verify pawn movement and marker behavior in the running app
- `/grill-me` — if questions arise about renderer integration (pawn/marker visuals) or level design tooling
- `/code-review` — after implementing simulation logic, especially `navigation.go` edge transitions (easy to get wrong)
- `caveman:cavecrew-builder` — for single-file edits (e.g. updating cube.go, writing loader.go)
- `caveman:cavecrew-investigator` — to locate where to hook simulation into the render/game loop
