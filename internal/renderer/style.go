package renderer

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Style struct {
	BackgroundColor rl.Color
	FaceColor       rl.Color
	GridColor       rl.Color
	HoverColor      rl.Color
}

func DefaultStyle() Style {
	return Style{
		BackgroundColor: rl.Black,
		FaceColor:       rl.White,
		GridColor:       rl.DarkGray,
		HoverColor:      rl.Orange,
	}
}

func LoadStyle(path string) (Style, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Style{}, fmt.Errorf("LoadStyle: %w", err)
	}

	var raw struct {
		BackgroundColor string `json:"background_color"`
		FaceColor       string `json:"face_color"`
		GridColor       string `json:"grid_color"`
		HoverColor      string `json:"hover_color"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return Style{}, fmt.Errorf("LoadStyle: %w", err)
	}

	bg, err := parseHexColor(raw.BackgroundColor)
	if err != nil {
		return Style{}, fmt.Errorf("LoadStyle background_color: %w", err)
	}
	fc, err := parseHexColor(raw.FaceColor)
	if err != nil {
		return Style{}, fmt.Errorf("LoadStyle face_color: %w", err)
	}
	gc, err := parseHexColor(raw.GridColor)
	if err != nil {
		return Style{}, fmt.Errorf("LoadStyle grid_color: %w", err)
	}
	hc, err := parseHexColor(raw.HoverColor)
	if err != nil {
		return Style{}, fmt.Errorf("LoadStyle hover_color: %w", err)
	}

	return Style{BackgroundColor: bg, FaceColor: fc, GridColor: gc, HoverColor: hc}, nil
}

func parseHexColor(s string) (rl.Color, error) {
	s = strings.TrimPrefix(s, "#")
	switch len(s) {
	case 6:
		v, err := strconv.ParseUint(s, 16, 32)
		if err != nil {
			return rl.Color{}, fmt.Errorf("invalid hex color %q: %w", s, err)
		}
		return rl.NewColor(uint8(v>>16), uint8(v>>8), uint8(v), 255), nil
	case 8:
		v, err := strconv.ParseUint(s, 16, 64)
		if err != nil {
			return rl.Color{}, fmt.Errorf("invalid hex color %q: %w", s, err)
		}
		return rl.NewColor(uint8(v>>24), uint8(v>>16), uint8(v>>8), uint8(v)), nil
	default:
		return rl.Color{}, fmt.Errorf("invalid hex color %q: expected #RRGGBB or #RRGGBBAA", s)
	}
}
