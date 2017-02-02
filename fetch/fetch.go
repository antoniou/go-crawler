package fetch

import (
	"fmt"
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

	// Retrieve Worker that manages Fetcher Service
	Worker() Worker
}

// Message is a struct used to pass results of a Fetch
// request back to the requester. It includes
// Request: The original request (for tracking)
// Response
// Error in case request could not finish successfully
type Message struct {
	Request  *url.URL
	Response *http.Response
	Error    error
}

// RequestQueue is used for incoming
// requests to the fetcher
type RequestQueue chan url.URL

// ResponseQueue queue is used for outgoing
// responses from the Fetcher
type ResponseQueue chan *Message

// AsyncHTTPFetcher implements Fetcher
type AsyncHTTPFetcher struct {
	//AsyncHTTPFetcher is an Asynchronous Worker
	*AsyncWorker

	requestQueue  *RequestQueue
	responseQueue *ResponseQueue

	client HTPPClient
}

// NewAsyncHTTPFetcher is a constructor for a
// AsyncHTTPFetcher. It does not start the
// Fetcher, which should be done by using the
// Run method
func NewAsyncHTTPFetcher() *AsyncHTTPFetcher {
	reqQueue := make(RequestQueue)
	resQueue := make(ResponseQueue)
	a := &AsyncHTTPFetcher{
		AsyncWorker: NewAsyncWorker("Fetcher"),

		client:        &http.Client{},
		requestQueue:  &reqQueue,
		responseQueue: &resQueue,
	}
	a.AsyncWorker.RunFunc = a.Run

	return a
}

// Fetch places a request for a URL into the requestQueue
// Returns nil on success and an error in case the url
// is not valid
func (a *AsyncHTTPFetcher) Fetch(url *url.URL) error {
	if err := a.validate(url); err != nil {
		return err
	}
	fmt.Printf("Fetcher: Adding URL %v to request queue\n", url)
	*a.requestQueue <- *url
	return nil
}

// Retrieve blocks waiting for responses to
// requests previously created with Fetch
func (a *AsyncHTTPFetcher) Retrieve() (Response *Message, err error) {
	Response = <-*a.responseQueue
	fmt.Printf("Fetcher: Passing result to Parser %v\n", *Response)
	return Response, Response.Error
}

// FIXME - This might not be needed
func (a *AsyncHTTPFetcher) get(url url.URL) (Response *http.Response, err error) {
	return a.client.Get(url.String())
}

// Worker Returns the embedded AsyncWorker struct
// which is used to Run and Stop the fetcher worker
func (a *AsyncHTTPFetcher) Worker() Worker {
	return a.AsyncWorker
}

// Run starts a loop that waits for requests
// or the quit signal. Run will be interrupted
// once the Stop method is used
func (a *AsyncHTTPFetcher) Run() error {
	a.AsyncWorker.SetState(RUNNING)
	for {
		a.AsyncWorker.SetState(WAITING)
		select {

		// A request is received
		case req := <-*a.requestQueue:
			a.AsyncWorker.SetState(RUNNING)
			res, err := a.get(req)
			*a.responseQueue <- &Message{
				Request:  &req,
				Response: res,
				Error:    err,
			}

		// A quit has been received, Stop has been invoked
		case <-a.AsyncWorker.quit:
			a.AsyncWorker.SetState(STOPPED)
			return nil
		default:
		}
	}
}

func (a *AsyncHTTPFetcher) validate(uri *url.URL) error {
	if uri.Scheme != "http" && uri.Scheme != "https" {
		return fmt.Errorf("Unsupported uri scheme %s", uri.Scheme)
	}

	return nil
}
