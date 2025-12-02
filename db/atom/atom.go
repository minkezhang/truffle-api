// Package atom is a collection of discrete data sources used to represent
// different data types.
//
// Each data type has a specific associated atom.
package atom

type AtomType int
type ClientAPI int

const (
	AtomTypeNone AtomType = iota
	AtomTypeTV
)

const (
	ClientAPIVirtual ClientAPI = iota
	ClientAPIBene
)

var (
	_ A[*Base] = &Base{}
)

type A[T any] interface {
	Type() AtomType
	API() ClientAPI
	ID() string
	GetBase() *Base
	Merge(T) T
}

type T struct {
	Title        string
	Localization string
}

type O struct {
	API        ClientAPI
	ID         string
	Titles     []T
	PreviewURL string
	Score      int
}

func New(o O) (*Base, error) {
	return &Base{
		api:        o.API,
		id:         o.ID,
		Titles:     o.Titles,
		PreviewURL: o.PreviewURL,
		Score:      o.Score,
	}, nil
}

type Base struct {
	api ClientAPI // Read-only
	id  string    // Read-only

	Titles     []T
	PreviewURL string
	Score      int
}

func (t *Base) API() ClientAPI { return t.api }
func (t *Base) ID() string     { return t.id }
func (t *Base) Type() AtomType { return AtomTypeNone }
func (t *Base) GetBase() *Base { return t }

func (t *Base) Merge(other *Base) *Base {
	if t == nil {
		t = &Base{}
	}
	return &Base{
		api: t.api,
		id:  t.id,
		Titles: append(
			append([]T{}, t.Titles...),
			other.Titles...),
		PreviewURL: other.PreviewURL,
		Score:      other.Score,
	}
}
