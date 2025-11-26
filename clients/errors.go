package client

import (
	"fmt"
	"net/http"

	apb "github.com/minkezhang/bene-api/proto/go/api"
)

type UnsupportedAPI struct {
	API apb.API
}

func (e UnsupportedAPI) Error() string {
	var api string
	if _, ok := apb.API_name[int32(e.API)]; !ok {
		api = apb.API_name[int32(apb.API_API_UNKNOWN)]
	} else {
		api = apb.API_name[int32(e.API)]
	}
	return fmt.Sprintf("Unsupported API: %s", api)
}

func (e UnsupportedAPI) Status() int {
	return http.StatusBadRequest
}
