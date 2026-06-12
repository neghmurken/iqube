# iQube — Game Specification

## Overview

iQube is a puzzle game played on a 3D cube. Each face of the cube is divided into a square grid of cells. The player must guide a pawn from a starting cell to a goal cell by placing directional markers on the grid before the simulation runs.

## The Cube

The cube has 6 faces: Front, Back, Left, Right, Top, Bottom. Each face is an n×n grid of cells.

Cells have one of three kinds:

| Kind | Color | Effect |
|------|-------|--------|
| **Normal** | face color | Pawn can cross freely. |
| **Void** | dark red | Impassable. Pawn steps onto it, then the simulation stops and placement resumes. The pawn returns to the start cell; placed markers are preserved. |
| **Blocked** | dark gray | Impassable. Pawn stays on the previous cell and retries every tick. Simulation keeps running. Movement resumes if the cell kind changes. |

## The Pawn

The pawn starts each level on a designated **start cell**, facing a fixed **initial direction**. It moves one cell at a time in a straight line. The pawn cannot choose its own path — it simply follows its current direction until something changes it.

When the pawn reaches the edge of a face, it crosses onto the adjacent face. The entry point and direction on the new face are determined by the geometry of the cube — the path remains continuous across the surface.

## Objective

The player wins the level when the pawn reaches the **goal cell**.

## Markers

Markers are the player's only tool. They are placed on empty cells and instruct the pawn to change direction when it steps onto them. Three types exist:

| Marker | Effect |
|--------|--------|
| **Turn Left** | Rotates the pawn 90° to its left |
| **Turn Right** | Rotates the pawn 90° to its right |
| **Turn Around** | Reverses the pawn's direction (180°) |

Direction changes are relative to the pawn's current direction, not to the cube's faces.

Each level provides a fixed **inventory** of markers — a limited count for each type. Markers cannot be placed if the corresponding inventory is empty.

Markers persist across the simulation: the pawn does not consume a marker when it crosses it. The same marker can redirect the pawn multiple times if the pawn loops back over it.

Only one marker can occupy a cell at a time.

Markers cannot be placed on the start cell, the goal cell, or non-normal cells (void or blocked).

## Placement Phase

Before starting the simulation, the player places markers on the grid. Hovering over a cell highlights it; selecting a marker type (1, 2, 3) places it on the hovered cell.

**Placing on an occupied cell:**
- Same type as the existing marker → removes it, returning it to the inventory.
- Different type → replaces it; the old marker returns to the inventory, the new one is spent.

The player can **reset** at any time during placement: all placed markers are removed and returned to the inventory, and the pawn returns to the start position.

## Simulation Phase

Once the player starts the simulation, the pawn begins moving. Markers can no longer be placed or removed.

Each step:
1. The pawn moves one cell in its current direction.
2. If the destination cell is **void**, the pawn enters it (the move is animated), then the simulation stops and placement resumes. The pawn returns to the start cell; placed markers are preserved.
3. If the destination cell is **blocked**, the pawn stays on the previous cell and retries on the next tick. Simulation continues.
4. If the destination cell is the **goal**, the level is won.
5. If the destination cell has a **marker**, the pawn's direction changes according to the marker's instruction. The change takes effect on the next step.

When crossing a face edge, the pawn's direction may also change to remain consistent with the cube's geometry.

The player can **stop** the simulation at any time, returning to placement with markers intact and the pawn back at the start.

## Winning and Progression

When the pawn reaches the goal cell, the level is won. The player proceeds to the next level.

Levels are played in sequence. After the last level, the sequence loops back to the first.

## Failure and Iteration

There is no explicit failure state. If the pawn is blocked or the player stops the simulation, placement resumes and the configuration can be adjusted. The player can iterate freely until the level is solved.

Note: the game does not currently detect infinite loops. If the pawn loops indefinitely without reaching the goal or a void cell, the simulation runs forever until the player stops it manually.
