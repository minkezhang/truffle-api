# truffle-api
Rewrite of https://github.com/minkezhang/truffle.

## Example

```golang
import (
	"github.com/minkezhang/truffle-api/client/option"
	"github.com/minkezhang/truffle-api/data/node"
	"github.com/minkezhang/truffle-api/data/source"
    "github.com/minkezhang/truffle-api/db"

	cpb "github.com/minkezhang/truffle-api/proto/go/config"
	dpb "github.com/minkezhang/truffle-api/proto/go/data"
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

func main() {
    db := db.New(
      context.Background(),
      &cpb.Config{ ... },   // Read from file
      &dpb.Database{ ... }, // Read from file
    )

    results, _ := db.Search(
      context.Background(),
      "frieren",
      map[epb.SourceAPI][]option.O{
        // Do an online search.
        epb.SourceAPI_SOURCE_API_MAL: []option.O{
          option.Remote(true),
          option.NSFW(false),
        },
        epb.SourceAPI_SOURCE_API_TRUFFLE: []option.O{},
      },
    )

    // Add and link to a new Truffle node. Multiple sources may be linked to the
    // same node as long as they have match source.S.SourceType.
    db.Put(context.Background(), source.Make(results[0]))

    // Get a specific node from a specific API.
    s, _ := db.Get(
      context.Background(),
      results[0].Header(),  // Includes SourceAPI, e.g. MAL.
      option.Remote(false), // Get source from cache.
    )

    // Get logical grouping of this source and related nodes, e.g. custom
    // Truffle data.
    n, _ := db.GetNode(
      context.Background(),
      node.Make(&dpb.Node{
        Header: &dpb.NodeHeader{
          Id:   s.NodeID(),
          Type: s.Header().Type(),
        },
      }).Header(),
    )
}
```
