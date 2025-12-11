# truffle-api
Rewrite of https://github.com/minkezhang/truffle.

## Adding a New Atom Type

1. Add new metadata message to `/proto/metadata.proto`
1. Add message as a field to `/proto/atom.proto`
1. Create a package in `/db/atom/metadata/${TYPE}/`
1. Implement the `metadata.G` and `metadata.M` interfaces for this new type
1. Update `atom.MergeMetadata`, `atom.Load`, and `atom.Save`

## Example

```golang
import (
    "github.com/minkezhang/truffle-api/client/mal"
    "github.com/minkezhang/truffle-api/db"
    "github.com/minkezhang/truffle-api/db/node"
)

func main() {
    d, err := db.New(context.Background(), db.O{  // TODO: Implement db.Load()
                                                  // and db.Save()
        Data:    []*node.N{},
        Clients: []client.C{
            mal.New(mal.O{
                ClientID:         "...",  // Apply for an API key via
                                          // https://help.myanimelist.net/hc/en-us/articles/900003108823-API
                PopularityCutoff: 5000,
                MaxResults:       10,
                NSFW:             true,
            }),
            ...  // TODO: Implement e.g. Steam, OMDB, Spotify clients etc.
        },
    })
    
}
```
