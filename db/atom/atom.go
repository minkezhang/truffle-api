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

	ClientAPINone ClientAPI = iota
	ClientTypeVirtual
	ClientAPIBene
)

var (
	_ A[*Base] = &Base{}
)

type A[T any] interface {
	Type() AtomType
	API() ClientAPI
	ID() string
	Merge(T) T
}

type T struct {
	Title        string
	Localization string
}

type H struct {
	API ClientAPI
	ID  string
}

type O struct {
	Titles     []T
	PreviewURL string
	Score      int
}

func New(o O) *Base {
	return &Base{
		Titles:     o.Titles,
		PreviewURL: o.PreviewURL,
		Score:      o.Score,
	}
}

type Base struct {
	api ClientAPI // Read-only
	id  string    // Read-only

	Titles     []T
	PreviewURL string
	Score      int
}

func (a *Base) API() ClientAPI { return a.api }
func (a *Base) ID() string     { return a.id }
func (a *Base) Type() AtomType { return AtomTypeNone }
func (a *Base) Merge(other *Base) *Base {
	res := &Base{
		api: a.api,
		id:  a.id,
		Titles: append(
			append([]T{}, a.Titles...),
			other.Titles...),
		PreviewURL: other.PreviewURL,
		Score:      other.Score,
	}
	return res
}
func (a *Base) WithHeader(h H) *Base {
	a.api = h.API
	a.id = h.ID
	return a
}
