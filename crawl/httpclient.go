package crawl

import "net/http"

// HTPPClient interface that wraps around the http.Client
// struct and can be replaced by any other client implementation
type HTPPClient interface {
	// At the moment, response is of type http.Response which locks
	// in implementation!
	Get(url string) (resp *http.Response, err error)
}

// Responder interface encapsulates the needed
// behaviour of http.Response. Rather than using
// the http.Response struct, the consumers use this
// interface instead
// type Responder interface {
//
// 	// Body of the Response
// 	BodyB() *[]byte
//
// 	// Status of the Response
// 	Status() string
// }

// type Response struct {
// 	*http.Response
// }
//
// func (r *Response) BodyB() *[]byte {
//
// 	r.Body.Read(p)
//
// }
