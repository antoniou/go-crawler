package crawl

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/antoniou/go-crawler/util"
	"github.com/goware/urlx"
)

// Fetcher is an Asynchronous interface
type Fetcher interface {
	// Fetch provides work to the Fetcher, in the
	// form of a URL to process
	Fetch(url *url.URL) error

	// ResponseChannel is a Getter returning
	// the Fetcher's Channel  that consumers
	// should be receiving results from
	ResponseChannel() (responseQueue *FetchResponseQueue)

	// Retrieve Worker that manages Fetcher Service
	Worker() Worker
}

// FetchMessage is a struct used to pass results of a Fetch
// request back to the requester. It includes
// Request: The original request (for tracking)
// Response
// Error in case request could not finish successfully
type FetchMessage struct {
	Request  *url.URL
	Response *http.Response
	Error    error
}

// RequestQueue is used for incoming
// requests to the fetcher
type RequestQueue chan url.URL

// FetchResponseQueue queue is used for outgoing
// responses from the Fetcher
type FetchResponseQueue chan *FetchMessage

// AsyncHTTPFetcher implements Fetcher
type AsyncHTTPFetcher struct {
	//AsyncHTTPFetcher is an Asynchronous Worker
	*AsyncWorker

	requestQueue  *RequestQueue
	responseQueue *FetchResponseQueue

	client HTPPClient
}

// NewAsyncHTTPFetcher is a constructor for a
// AsyncHTTPFetcher. It does not start the
// Fetcher, which should be done by using the
// Run method
func NewAsyncHTTPFetcher() *AsyncHTTPFetcher {
	reqQueue := make(RequestQueue)
	resQueue := make(FetchResponseQueue)
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
	if a.AsyncWorker.State() == STOPPED {
		return fmt.Errorf("%s is in state stopped", a.AsyncWorker.Type())
	}
	if err := a.validate(url); err != nil {
		return err
	}
	normURL, _ := urlx.Normalize(url)
	util.Printf("Fetcher: Adding URL %v to request queue\n", normURL)
	normURLs, _ := urlx.Parse(normURL)
	*a.requestQueue <- *normURLs
	return nil
}

// ResponseChannel is a Getter returning the Fetcher's Channel  that consumers
// should be receiving results from
func (a *AsyncHTTPFetcher) ResponseChannel() (responseQueue *FetchResponseQueue) {
	return a.responseQueue
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
			res, err := a.client.Get(req.String())
			*a.responseQueue <- &FetchMessage{
				Request:  &req,
				Response: res,
				Error:    err,
			}

			// A quit has been received, Stop has been invoked
		case <-a.AsyncWorker.Quit:
			a.Worker().SetState(STOPPED)
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
