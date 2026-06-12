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

## Commands

- `make build` build the app
- `make run` build and run the app
