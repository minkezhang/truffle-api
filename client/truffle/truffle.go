package truffle

import (
	"context"
	"fmt"
	"regexp"

	"github.com/minkezhang/truffle-api/client/option"
	"github.com/minkezhang/truffle-api/data/source"
	"github.com/minkezhang/truffle-api/util/generator"
	"github.com/minkezhang/truffle-api/util/slice"

	cpb "github.com/minkezhang/truffle-api/proto/go/config"
	dpb "github.com/minkezhang/truffle-api/proto/go/data"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

func New(pb *cpb.Truffle, data []*dpb.Source) *C {
	ids := slice.Apply(data, func(v *dpb.Source) string {
		return v.GetHeader().GetId()
	})

	sources := map[string]source.S{}
	for _, s := range data {
		sources[s.GetHeader().GetId()] = source.Make(s)
	}

	return &C{
		sources: sources,
		generator: generator.New(generator.O{
			IDs: ids,
			N:   16,
		}),
	}
}

type C struct {
	sources   map[string]source.S
	generator *generator.G
}

func (c *C) Put(ctx context.Context, s source.S) (source.H, error) {
	pb := s.PB()
	if pb.GetHeader().GetId() == "" {
		pb.Header = &dpb.SourceHeader{
			Id:   c.generator.Generate(),
			Type: pb.GetHeader().GetType(),
			Api:  epb.SourceAPI_SOURCE_API_MAL,
		}
	}
	s = source.Make(pb)
	c.sources[s.Header().ID()] = s
	return s.Header(), nil
}

func (c *C) Delete(ctx context.Context, header source.H) error {
	delete(c.sources, header.ID())
	return nil
}

func (c *C) Get(ctx context.Context, header source.H) (source.S, error) {
	return c.sources[header.ID()], nil
}

func (c *C) Search(ctx context.Context, query string, opts ...option.O) ([]source.S, error) {
	pattern, err := regexp.Compile(fmt.Sprintf("(?i)%v", query))
	if err != nil {
		return nil, err
	}
	var results []source.S
	for _, s := range c.sources {
		for _, t := range s.Titles() {
			if pattern.MatchString(t.Title()) {
				results = append(results, s)
				break
			}
		}
	}

	return results, nil
}
