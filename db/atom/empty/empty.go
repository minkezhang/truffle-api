// Package empty is used to initialize an `atom.A` without any specific
// auxiliary data.
//
// Example:
//
//	atom.New(atom.O{
//		...
//		Aux: empty.A{},
//	})
package empty

import (
	"github.com/minkezhang/bene-api/db/atom"
	"github.com/minkezhang/bene-api/db/enums"
)

var (
	_ atom.Aux = A{}
)

type A struct{}

func (a A) AtomType() enums.AtomType  { return enums.AtomTypeNone }
func (a A) Equal(o atom.Aux) bool     { return true }
func (a A) Copy() atom.Aux            { return A{} }
func (a A) Merge(o atom.Aux) atom.Aux { return o.Copy() }
