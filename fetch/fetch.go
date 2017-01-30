package fetch

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Fetcher is an Asynchronous interface
type Fetcher interface {
	// Fetch provides work to the Fetcher, in the
	// form of a URL to process
	Fetch(url *url.URL) error

	// Retrieve provides results back from the Fetcher
	// in the form of a Response
	Retrieve() (Response *Message, err error)

	// Run starts the Fetcher
	Run() error
}

// HttpClient interface is
type HttpClient interface {
	// Get
	Get(url string) (resp *http.Response, err error)
}

func NewAsyncHttpFetcher() *AsyncHttpFetcher {
	reqQueue := make(RequestQueue)
	resQueue := make(ResponseQueue)
	a := &AsyncHttpFetcher{
		client:        &http.Client{},
		requestQueue:  &reqQueue,
		responseQueue: &resQueue,
	}

	return a
}

type Message struct {
	Request  *url.URL
	Response io.ReadCloser
}

type RequestQueue chan url.URL
type ResponseQueue chan *Message

// AsyncHttpFetcher implements Fetcher
type AsyncHttpFetcher struct {
	requestQueue  *RequestQueue
	responseQueue *ResponseQueue

	client HttpClient
}

func (a *AsyncHttpFetcher) Fetch(url *url.URL) error {
	fmt.Printf("Fetcher: Adding URL %v to request queue\n", url)
	*a.requestQueue <- *url
	return nil
}

func (a *AsyncHttpFetcher) Retrieve() (Response *Message, err error) {
	Response = <-*a.responseQueue
	fmt.Printf("Fetcher: Passing result to Parser\n")
	return Response, nil
}

// FIXME - This might not be needed
func (a *AsyncHttpFetcher) get(url url.URL) (Response io.ReadCloser, err error) {
	res, err := a.client.Get(url.String())
	return res.Body, err
}

func (a *AsyncHttpFetcher) Run() error {
	for {
		select {
		case req := <-*a.requestQueue:
			fmt.Printf("Fetcher: Going to fetch url %s\n", req.String())
			res, _ := a.get(req)
			// fmt.Printf("Fetcher: Waiting for result of get to %s\n", req.String())
			*a.responseQueue <- &Message{
				Request:  &req,
				Response: res,
			}
			fmt.Printf("Fetcher: Got result of get to %s\n", req.String())
		default:

		}
	}
}
