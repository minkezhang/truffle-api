package enums

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
