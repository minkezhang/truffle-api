package atom

import (
	"fmt"

	apb "github.com/minkezhang/bene-api/proto/go/api"
	dpb "github.com/minkezhang/bene-api/proto/go/data"
)

type T struct {
	Title string
	Localization string
}

type A struct {
	nodeType dpb.NodeType
	nodeID   string

	api        apb.API
	id         string
	isDirty    bool

	titles []T
	previewURL string
	score      int

	season int
	genres []string
	showrunners []string
	isAnimated bool
	directors []string
	writers []string
	cinematography []string
	composers []string
	starring []string
	animationStudios []string
	creator  string
	videoURL string
}

func (b *A) NodeType() dpb.NodeType { return b.nodeType }
func (b *A) NodeID() string         { return b.nodeID }
func (b *A) API() apb.API           { return b.api }
func (b *A) ID() string             { return b.id }
func (b *A) IsDirty() bool          { return b.isDirty }
func (b *A) PreviewURL() string     { return b.previewURL }
func (b *A) Score() int             { return b.score }

func (b *A) Update(other A) error {
	if b.API() != other.API() || b.ID() != other.ID() {
		return fmt.Errorf("cannot fulfill a PUT request on mismatching atoms: (%v:%v != %v:%v)", b.API(), b.ID(), other.API(), other.ID())
	}
	return nil
}
