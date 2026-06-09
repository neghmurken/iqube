# iqube

iQube is a puzzle game written in Go and Raylib.

## Description

The goal of the game is to move a piece across a cube, each face of which is divided into a grid of cells.
At the start of the game, the piece is on a starting cell and must move to a finishing cell, which is the objective of the level.
Cells can be of several types: empty or filled. The pawn can move through empty cells, but not through filled cells.

The player must then place markers on the empty cells that give instructions to the pawn to change directions.

These markers are limited for each level and can be of several types:
 - turn left
 - turn right
 - turn around

The pawn can only move in one direction at the start.
As soon as the pawn passes over a cell containing a marker, it executes the marker’s instruction and changes direction accordingly.

The player must first place all or some of their markers and then start the simulation for the pawn to move. Once the simulation is running, it is no longer possible to place or remove markers.
If the pawn reaches the goal, the level is won. If a pawn encounters a cell that is full, the pawn is blocked and the simulation stops, allowing the player to modify their marker configuration.
If the pawn reaches the edge of one side of the cube, the pawn continues along its path on the adjacent side.
