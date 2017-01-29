package fetch

import (
	"fmt"
	"io"
	"net/http"
)

// Fetcher is an Asynchronous interface
type Fetcher interface {
	// Fetch provides work to the Fetcher, in the
	// form of a URL to process
	Fetch(url string) error

	// Retrieve provides results back from the Fetcher
	// in the form of a Response
	Retrieve() (Response io.ReadCloser, err error)
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
	a.Run()

	return a
}

type RequestQueue chan string
type ResponseQueue chan io.ReadCloser

// AsyncHttpFetcher implements Fetcher
type AsyncHttpFetcher struct {
	requestQueue  *RequestQueue
	responseQueue *ResponseQueue

	client HttpClient
}

func (a *AsyncHttpFetcher) Fetch(url string) error {
	fmt.Printf("Fetcher: Adding URL %s to request queue\n", url)
	*a.requestQueue <- url
	return nil
}

func (a *AsyncHttpFetcher) Retrieve() (Response io.ReadCloser, err error) {
	fmt.Printf("Fetcher: Waiting for response...")
	Response = <-*a.responseQueue
	fmt.Printf("Done\n")
	return Response, nil
}

// FIXME - This might not be needed
func (a *AsyncHttpFetcher) get(url string) (Response io.ReadCloser, err error) {
	res, err := a.client.Get(url)
	return res.Body, err
}

func (a *AsyncHttpFetcher) Run() {
	go a.background()
}

func (a *AsyncHttpFetcher) background() {
	for {
		select {
		case req := <-*a.requestQueue:
			res, _ := a.get(req)
			// sres := make([]byte, 1000000)
			// res.Read(sres)
			*a.responseQueue <- res
		default:

		}
	}
}
