package fetch

import "net/http"

// HTPPClient interface that wraps around the http.Client
// struct and can be replaced by any other client implementation
type HTPPClient interface {
	// At the moment, response is of type http.Response which locks
	// in implementation!
	Get(url string) (resp *http.Response, err error)
}
