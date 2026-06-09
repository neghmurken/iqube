# iqube

## General behavior

Ask questions if necessary, don't take decisions all by yourself but propose solutions and architecture.

## Dependencies 
 - Go
 - Raylib (docs: https://www.raylib.com/cheatsheet/cheatsheet.html)

## Architecture

```
main.go                          Entry point: window init, game loop (BeginDrawing / EndDrawing)
Makefile                         build / run recipes

assets/
  theme.json                     Scene colors: background_color, face_color, grid_color (#RRGGBB or #RRGGBBAA)

internal/
  game/
    game.go                      Game struct — owns Cube + Renderer, wires Update/Draw
  model/
    cube.go                      Cube{GridSize, Faces[6][][]Cell}, face index constants (FacePosX…FaceNegZ)
  renderer/
    renderer.go                  Renderer: orbit camera (spherical coords), right-drag rotation,
                                 scroll-wheel zoom, 2-pass 3D draw (filled faces then grid lines)
    style.go                     Style{BackgroundColor, FaceColor, GridColor},
                                 DefaultStyle(), LoadStyle(path), parseHexColor()
```

### Rendering pipeline

1. `game.Draw()` → `renderer.Draw(cube)`
2. `ClearBackground` with `style.BackgroundColor`
3. `BeginMode3D` (perspective camera, orbit around origin)
4. Pass 1 — front-facing faces only (`normal · cameraPos > 0`): fill each face with 2 triangles (`style.FaceColor`)
5. Pass 2 — same faces: draw N+1 grid lines in each axis, offset by `normal × lineOffset` to avoid z-fighting (`style.GridColor`)

### Camera

Stored as spherical coords `(radius, azimuth, altitude)`, synced to `rl.Camera3D.Position` each frame.
- Right-drag → azimuth / altitude
- Scroll wheel → radius, clamped to `[minRadius, maxRadius]`

## Specs

Follow the specs -> SPECS.md

# Commands

- `make build` build the app
- `make run` build and run the app