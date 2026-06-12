# iqube

## General behavior

Ask if unclear. Propose solutions/architecture, don't decide alone.

## Specs

General picture: @SPECS.md

## Dependencies 
 - Go
 - Raylib (docs: https://www.raylib.com/cheatsheet/cheatsheet.html)

## Architecture

```
main.go                          Entry point: window init, game loop (BeginDrawing / EndDrawing)
Makefile                         build / run recipes

assets/
  theme.json                     Scene colors, font, marker sprite paths (#RRGGBB or #RRGGBBAA)

internal/
  game/
    game.go                      Game struct — owns Cube + Renderer, wires Update/Draw
  model/
    cube.go                      Cube{GridSize, Faces[6][][]Cell}, face index constants (RG…BK)
    state.go                     GamePhase (Placement/Running/Won), State (game state snapshot)
  renderer/
    renderer.go                  Renderer struct, New/Unload/Update/Draw entry points;
                                 State embeds model.State + PawnPrev/PawnAnimT (animation)
    camera.go                    Spherical camera constants, handleCameraInput, syncCameraPosition
    face.go                      faceDesc, buildFaces, hitTest (ray → cell)
    draw3d.go                    drawCubeGrid, drawFace, drawCell, drawFaceGrid, drawFaceOrigin
    draw_entities.go             drawPawn, drawMarker, cellColor, cellWorldCenter
    hud.go                       drawHUD (2D overlay: phase, inventory, hints, hover label)
    style.go                     Style{Colors, FontPath, MarkerSpritePath, FontSize},
                                 DefaultStyle(), LoadStyle(path), parseHexColor()
```

### Glossary

Cube has 6 faces:
 - UP: top face
 - BT: bottom face
 - LF: left face
 - RG: right face
 - FR: front face
 - BK: back face

Each face: n*n cell grid. Origin = bottom-left cell. Cell named `{FaceName}({col},{row})`. Example: `BT(4,2)`

### Rendering pipeline

1. `game.Draw()` → `renderer.Draw(cube)`
2. `ClearBackground` with `style.BackgroundColor`
3. `BeginMode3D` (perspective camera, orbit around origin)
4. Pass 1 — front-facing faces only (`normal · cameraPos > 0`): fill each face with 2 triangles (`style.FaceColor`)
5. Pass 2 — same faces: draw N+1 grid lines per axis, offset by `normal × lineOffset` to avoid z-fighting (`style.GridColor`)

### Camera

Spherical coords `(radius, azimuth, altitude)`, synced to `rl.Camera3D.Position` each frame.
- Right-drag → azimuth / altitude
- Scroll wheel → radius, clamped to `[minRadius, maxRadius]`

## Commands

- `make build` build the app
- `make run` build and run the app

## Agent skills

### Issue tracker

Issues in GitHub Issues (`gh` CLI). See `docs/agents/issue-tracker.md`.

### Triage labels

Default canonical label strings — no custom overrides. See `docs/agents/triage-labels.md`.

### Domain docs

Single-context repo: one `CONTEXT.md` + `docs/adr/` at root. See `docs/agents/domain.md`.