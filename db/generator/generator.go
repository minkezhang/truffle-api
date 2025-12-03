package generator

import (
	"math/rand"
	"time"
)

const (
	l = 32
)

var (
	runes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
)

type O struct {
	IDs  []string
	Seed int64
	N    int
}

type G struct {
	ids  map[string]bool
	rand *rand.Rand
	n    int
}

func New(o O) *G {
	s := o.Seed
	if s == 0 {
		s = time.Now().UnixNano()
	}
	n := o.N
	if n == 0 {
		n = l
	}

	g := &G{
		ids:  map[string]bool{},
		rand: rand.New(rand.NewSource(s)),
		n:    n,
	}
	for _, id := range o.IDs {
		g.ids[id] = true
	}
	return g
}

func (g *G) Generate() string {
	id := ""
	buf := make([]rune, g.n)
	for id == "" || g.ids[id] {
		for i := range buf {
			buf[i] = runes[g.rand.Intn(len(runes))]
		}
		id = string(buf)
	}
	g.ids[id] = true
	return id
}
