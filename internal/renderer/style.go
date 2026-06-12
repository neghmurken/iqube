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
	Colors           map[string]rl.Color
	MarkerSpritePath string
	FontPath         string
	FontSize         int32
}

func (s Style) GetColor(name string) rl.Color {
	return s.Colors[name]
}

func DefaultStyle() Style {
	return Style{
		Colors: map[string]rl.Color{
			"background": rl.Black,
			"face":       rl.White,
			"grid":       rl.DarkGray,
			"hover":      rl.Orange,
			"filled":     rl.NewColor(60, 60, 60, 255),
			"start":      rl.NewColor(80, 200, 80, 200),
			"goal":       rl.NewColor(255, 215, 0, 220),
			"pawn":       rl.Blue,
		},
		MarkerSpritePath: "assets/markers.png",
		FontSize:         20,
	}
}

func LoadStyle(path string) (Style, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Style{}, fmt.Errorf("LoadStyle: %w", err)
	}

	var raw struct {
		Colors       map[string]string `json:"colors"`
		MarkerSprite string            `json:"marker_sprite"`
		Font         string            `json:"font"`
		FontSize     int               `json:"font_size"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return Style{}, fmt.Errorf("LoadStyle: %w", err)
	}

	def := DefaultStyle()

	colors := make(map[string]rl.Color, len(def.Colors))
	for k, v := range def.Colors {
		colors[k] = v
	}
	for k, hex := range raw.Colors {
		c, err := parseHexColor(hex)
		if err != nil {
			return Style{}, fmt.Errorf("LoadStyle colors.%s: %w", k, err)
		}
		colors[k] = c
	}

	markerSprite := def.MarkerSpritePath
	if raw.MarkerSprite != "" {
		markerSprite = raw.MarkerSprite
	}

	fontSize := def.FontSize
	if raw.FontSize > 0 {
		fontSize = int32(raw.FontSize)
	}

	return Style{
		Colors:           colors,
		MarkerSpritePath: markerSprite,
		FontPath:         raw.Font,
		FontSize:         fontSize,
	}, nil
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
