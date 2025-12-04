Client
----

This directory provides all the API hook logic to connect with each provider.

A client in this context is a piece of logic which can query a data source
backend and return sources of truth wrapped in a logical `atom.A`.

Note that a client is considered read-only -- `bene` do not expect to update
remote sources of truth.

```golang
type C interface {
    APIType() enums.ClientAPI  // e.g. enums.ClientAPIMAL
    Get(ctx context.Context, id string) (*atom.A, error)
	Query(ctx context.Context, q query.Q) ([]*atom.A, error)
}
```
