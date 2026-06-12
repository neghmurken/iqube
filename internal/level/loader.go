package level

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"gopkg.in/yaml.v3"

	"github.com/neghmurken/iqube/internal/model"
)

var faceNames = map[string]int{
	"UP": model.UP,
	"BT": model.BT,
	"LF": model.LF,
	"RG": model.RG,
	"FR": model.FR,
	"BK": model.BK,
}

func parseFace(s string) (int, error) {
	if idx, ok := faceNames[s]; ok {
		return idx, nil
	}
	return 0, fmt.Errorf("unknown face %q: valid values are UP BT LF RG FR BK", s)
}

// cellCoord unmarshals either an integer or "*" (wildcard = all indices).
type cellCoord struct {
	all bool
	val int
}

func (c *cellCoord) UnmarshalYAML(value *yaml.Node) error {
	if value.Value == "*" {
		c.all = true
		return nil
	}
	c.all = false
	return value.Decode(&c.val)
}

func (c cellCoord) indices(n int) []int {
	if c.all {
		idx := make([]int, n)
		for i := range idx {
			idx[i] = i
		}
		return idx
	}
	return []int{c.val}
}

type yamlPosition struct {
	Face string `yaml:"face"`
	Row  int    `yaml:"row"`
	Col  int    `yaml:"col"`
}

type yamlCell struct {
	Face string    `yaml:"face"`
	Row  cellCoord `yaml:"row"`
	Col  cellCoord `yaml:"col"`
	Kind string    `yaml:"kind"`
}

func (c *yamlCell) UnmarshalYAML(value *yaml.Node) error {
	c.Row = cellCoord{all: true}
	c.Col = cellCoord{all: true}
	type alias yamlCell
	return value.Decode((*alias)(c))
}

type yamlInventory struct {
	TurnLeft   int `yaml:"turn_left"`
	TurnRight  int `yaml:"turn_right"`
	TurnAround int `yaml:"turn_around"`
}

type yamlLevel struct {
	GridSize         int           `yaml:"grid_size"`
	Start            yamlPosition  `yaml:"start"`
	Goal             yamlPosition  `yaml:"goal"`
	InitialDirection string        `yaml:"initial_direction"`
	Inventory        yamlInventory `yaml:"inventory"`
	Cells            []yamlCell    `yaml:"cells"`
}

func LoadAll(dir string) ([]model.Level, error) {
	entries, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return nil, fmt.Errorf("level.LoadAll: %w", err)
	}
	sort.Strings(entries)

	levels := make([]model.Level, 0, len(entries))
	for _, path := range entries {
		lvl, err := loadFile(path)
		if err != nil {
			return nil, fmt.Errorf("level.LoadAll %s: %w", path, err)
		}
		levels = append(levels, lvl)
	}
	return levels, nil
}

func loadFile(path string) (model.Level, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return model.Level{}, err
	}
	var raw yamlLevel
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return model.Level{}, err
	}

	startFace, err := parseFace(raw.Start.Face)
	if err != nil {
		return model.Level{}, fmt.Errorf("start.face: %w", err)
	}
	goalFace, err := parseFace(raw.Goal.Face)
	if err != nil {
		return model.Level{}, fmt.Errorf("goal.face: %w", err)
	}

	n := raw.GridSize
	cells := make([]model.LevelCell, 0, len(raw.Cells))
	for _, c := range raw.Cells {
		cellFace, err := parseFace(c.Face)
		if err != nil {
			return model.Level{}, fmt.Errorf("cell.face: %w", err)
		}
		kind := model.CellEmpty
		if c.Kind == "filled" {
			kind = model.CellFilled
		}
		for _, row := range c.Row.indices(n) {
			for _, col := range c.Col.indices(n) {
				cells = append(cells, model.LevelCell{
					Position: model.Position{Face: cellFace, Row: row, Col: col},
					Kind:     kind,
				})
			}
		}
	}

	return model.Level{
		GridSize:         n,
		Start:            model.Position{Face: startFace, Row: raw.Start.Row, Col: raw.Start.Col},
		Goal:             model.Position{Face: goalFace, Row: raw.Goal.Row, Col: raw.Goal.Col},
		InitialDirection: model.ParseDirection(raw.InitialDirection),
		Inventory: model.Inventory{
			TurnLeft:   raw.Inventory.TurnLeft,
			TurnRight:  raw.Inventory.TurnRight,
			TurnAround: raw.Inventory.TurnAround,
		},
		Cells: cells,
	}, nil
}
