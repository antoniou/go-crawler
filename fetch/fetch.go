package fetch

import (
	"fmt"
	"net/http"
)

// Fetcher is an interface
type Fetcher interface {
	// Fetch provides work to the Fetcher
	Fetch(url string) error

	// Retrieve provides work back from the Fetcher
	Retrieve() (Response string, err error)
}

func NewAsyncHttpFetcher() *AsyncHttpFetcher {
	reqQueue := make(RequestQueue)
	resQueue := make(ResponseQueue)
	a := &AsyncHttpFetcher{
		requestQueue:  &reqQueue,
		responseQueue: &resQueue,
	}
	go a.Run()

	return a
}

type RequestQueue chan string
type ResponseQueue chan string

// AsyncHttpFetcher implements Fetcher
type AsyncHttpFetcher struct {
	requestQueue  *RequestQueue
	responseQueue *ResponseQueue
}

func (a *AsyncHttpFetcher) Fetch(url string) error {

	fmt.Printf("Adding URL %s to request queue\n", url)
	*a.requestQueue <- url
	return nil
}

func (a *AsyncHttpFetcher) Retrieve() (Response string, err error) {
	fmt.Printf("Waiting for response...")
	Response = <-*a.responseQueue
	fmt.Printf("Done\n")
	return Response, nil
}

func (a *AsyncHttpFetcher) fetch(url string) (Response string, err error) {
	res, err := http.Get(url)
	body := make([]byte, 100000)
	res.Body.Read(body)
	return string(body), err
}

func (a *AsyncHttpFetcher) Run() {
	for {
		select {
		case req := <-*a.requestQueue:
			fmt.Printf("Fetching Url %s\n", req)
			res, _ := a.fetch(req)
			*a.responseQueue <- res
		default:

		}
	}
}
