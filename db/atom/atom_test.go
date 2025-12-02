package atom

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/minkezhang/bene-api/db/enums"
)

var (
	_ Aux = &MockAuxTV{}
)

type MockAuxTV struct {
	producers []string
}

func (a *MockAuxTV) AtomType() enums.AtomType { return enums.AtomTypeTV }
func (a *MockAuxTV) Producers() []string      { return append([]string{}, a.producers...) }
func (a *MockAuxTV) Equal(o Aux) bool         { return reflect.DeepEqual(a, o) }

func (a *MockAuxTV) Copy() Aux {
	return &MockAuxTV{
		producers: append([]string{}, a.producers...),
	}
}

func (a *MockAuxTV) Merge(o Aux) Aux {
	if a.AtomType() != o.AtomType() {
		panic(fmt.Errorf("cannot merge mismatching atom types: %v != %v", a.AtomType(), o.AtomType()))
	}
	return &MockAuxTV{
		producers: append(
			append([]string{}, a.producers...),
			o.(*MockAuxTV).producers...,
		),
	}
}

func TestMerge(t *testing.T) {
	got := New(O{
		APIType: enums.ClientAPIBene,
		APIID:   "foo",
		Titles: []T{
			{Title: "Firefly"},
		},
		PreviewURL: "",
		Score:      91,
		AtomType:   enums.AtomTypeTV,
		Aux: &MockAuxTV{
			producers: []string{"Joss Whedon"},
		},
	}).Merge(New(O{
		APIType: enums.ClientAPIBene,
		APIID:   "foo",
		Titles: []T{
			{Title: "Firefly"},
		},
		PreviewURL: "overwrite",
		Score:      92,
		AtomType:   enums.AtomTypeTV,
		Aux: &MockAuxTV{
			producers: []string{"Tim Minear"},
		},
	}))

	want := New(O{
		APIType: enums.ClientAPIBene,
		APIID:   "foo",
		Titles: []T{
			{Title: "Firefly"}, // Remove duplicates
		},
		PreviewURL: "overwrite",
		Score:      92,
		AtomType:   enums.AtomTypeTV,
		Aux: &MockAuxTV{
			producers: []string{"Joss Whedon", "Tim Minear"},
		},
	})

	if diff := cmp.Diff(
		want,
		got,
		cmp.AllowUnexported(A{}, MockAuxTV{}),
	); diff != "" {
		t.Errorf("Merge() mismatch (-want +got):\n%s", diff)
	}
}
