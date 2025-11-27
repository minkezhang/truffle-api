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
	GetType() AtomType
	GetAPI() ClientAPI
	GetID() string
	Merge(T) (T, error)
}

type Base struct {
	Type AtomType  // Read-only
	API  ClientAPI // Read-only
	ID   string    // Read-only

	Titles []struct {
		Title        string
		Localization string
	}

	PreviewURL string
	Score      int
}

func (a *Base) GetType() AtomType { return a.Type }
func (a *Base) GetAPI() ClientAPI { return a.API }
func (a *Base) GetID() string     { return a.ID }
func (a *Base) Merge(other *Base) (*Base, error) {
	return &Base{
		Type:    a.Type,
		API:     a.API,
		ID:      a.ID,
		Titles: append(
			append([]struct {
				Title string
				Localization string
			}{},
				a.Titles...),
			other.Titles...),
		PreviewURL: other.PreviewURL,
		Score: other.Score,
	}, nil
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
