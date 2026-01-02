package option

import (
	epb "github.com/minkezhang/truffle-api/proto/go/enums"
)

type O interface {
	is_option()

	// IsSupported returns a list of supported API endpoints.
	//
	// This is for documentation only. Clients will need to implement
	// handling logic.
	IsSupported(v epb.SourceAPI) bool
}

// Remote is an option specifying if the request should read from remote
// databases, e.g. make an API call to MAL. If set to false, the client will
// search will only look at the local cache.
type Remote bool

func (o Remote) is_option()                       {}
func (o Remote) IsSupported(v epb.SourceAPI) bool { return true }

// NSFW specifies if a client.Search query should include NSFW results.
//
// Supported APIs: MAL
type NSFW bool

func (o NSFW) is_option() {}
func (o NSFW) IsSupported(v epb.SourceAPI) bool {
	return map[epb.SourceAPI]bool{
		epb.SourceAPI_SOURCE_API_MAL: true,
	}[v]
}
