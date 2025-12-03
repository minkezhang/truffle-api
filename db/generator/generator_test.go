package generator

import (
	"testing"
)

func TestGenerate(t *testing.T) {
	t.Run("Base", func(t *testing.T) {
		g := New(O{})
		id := g.Generate()
		if !g.ids[id] {
			t.Errorf("generated ID is not saved in generator")
		}
	})
	t.Run("Duplicate", func(t *testing.T) {
		magic := "qxDvSDXcNRoIzFFzZNvIbMatNzvBGEZ7"
		g := New(O{
			IDs:  []string{magic},
			Seed: 2025,
		})
		id := g.Generate()
		if id == magic {
			t.Errorf("Generate() generated an already generated ID: %v", id)
		}
	})
}
