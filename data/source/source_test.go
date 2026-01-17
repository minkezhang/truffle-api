package source

import (
	"testing"
)

func TestWithNodeID(t *testing.T) {
	t.Run("NilPB", func(t *testing.T) {
		if got, want := (S{}).WithNodeID("foo").NodeID(), "foo"; got != want {
			t.Errorf("NodeID() = %v, want = %v", got, want)
		}
	})

}
