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

type Base struct {
	atomType AtomType  // Read-only
	api      ClientAPI // Read-only
	id       string    // Read-only

	Titles []struct {
		Title        string
		Localization string
	}

	PreviewURL string
	Score      int
}

func (a *Base) Type() AtomType { return a.atomType }
func (a *Base) API() ClientAPI { return a.api }
func (a *Base) ID() string     { return a.id }
func (a *Base) Merge(other *Base) *Base {
	res := &Base{
		atomType: a.atomType,
		api:      a.api,
		id:       a.id,
		Titles: append(
			append([]struct {
				Title        string
				Localization string
			}{},
				a.Titles...),
			other.Titles...),
		PreviewURL: other.PreviewURL,
		Score:      other.Score,
	}
	return res
}
func (a *Base) WithHeader(o struct {
	Type AtomType
	API  ClientAPI
	ID   string
}) {
	a.atomType = o.Type
	a.api = o.API
	a.id = o.ID
}

type TV struct {
	*Base

	Season           int
	IsAnimated       bool
	Genres           []string
	Showrunners      []string
	Directors        []string
	Writers          []string
	Cinematography   []string
	Composers        []string
	Starring         []string
	AnimationStudios []string
}
