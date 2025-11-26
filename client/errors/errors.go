package client

import (
	"fmt"
	"net/http"

	apb "github.com/minkezhang/bene-api/proto/go/api"
)

type ErrorOK struct{}

func (e ErrorOK) Error() string { return "" }
func (e ErrorOK) Status() int   { return http.StatusOK }

type ErrorUnsupportedAPI struct {
	API apb.API
}

func (e ErrorUnsupportedAPI) Error() string {
	var api string
	if _, ok := apb.API_name[int32(e.API)]; !ok {
		api = apb.API_name[int32(apb.API_API_UNKNOWN)]
	} else {
		api = apb.API_name[int32(e.API)]
	}
	return fmt.Sprintf("Unsupported API: %s", api)
}

func (e ErrorUnsupportedAPI) Status() int {
	return http.StatusBadRequest
}
