---
name: iqube-feature
description: Add a new feature to the iQube game. Gathers requirements, clarifies ambiguities, proposes an implementation plan, implements after user approval, then updates SPECS.md. Use when user says "add a feature", "I want to add X", "new feature", "let's implement X" in the context of the iqube game.
---

# iQube Feature Addition

## Process

1. **Understand the request** — ask what feature the user wants. If vague, ask 2–3 targeted questions to pin down scope and behavior before proposing anything.

2. **Clarify edge cases** — for each ambiguous interaction, ask explicitly. Common areas:
   - Does it affect simulation, placement, or both phases?
   - How does it interact with existing cell kinds (normal / void / blocked)?
   - How does it interact with markers and face transitions?
   - Does it need new YAML fields in level files?
   - Does it need new colors / assets in theme.json?

3. **Propose a plan** — list every file that changes and what changes in each. Format:
   ```
   src/model/…     — new types or constants
   src/simulation/ — step/advance logic
   src/renderer/   — draw + style + HUD
   src/level/      — YAML loader
   assets/         — theme.json, sprites
   SPECS.md        — documentation
   ```
   Wait for explicit user approval before touching any file.

4. **Implement** — follow CLAUDE.md conventions:
   - Go + Raylib (see cheatsheet link in CLAUDE.md)
   - Self-documenting code, no unnecessary comments
   - `make build` must pass after every file touched

5. **Update SPECS.md** — add or edit the relevant section. Keep tone consistent with existing sections (present tense, player-facing language).

## Key domain facts (read SPECS.md for full spec)

- Cube: 6 faces (UP/BT/LF/RG/FR/BK), each n×n grid
- Cell kinds: normal (passable), void (pawn enters then sim stops), blocked (pawn waits, sim continues)
- Markers: turn-left / turn-right / turn-around; placed in placement phase only; persist across simulation
- Pawn moves 1 cell/tick; direction changes take effect on next step
- Face transitions: geometry-driven, defined in `src/model/navigation.go`
- Renderer state: `PawnPrev` + `PawnAnimT` drive lerp animation between ticks

## Checklist before marking done

- [ ] `make build` passes
- [ ] New YAML fields documented with example in plan (if any)
- [ ] theme.json updated for new colors (if any)
- [ ] SPECS.md updated
