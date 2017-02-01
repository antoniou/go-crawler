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

	// Retrieve Worker that manages Fetcher Service
	Worker() Worker
}

type Message struct {
	Request  *url.URL
	Response io.ReadCloser
}

type RequestQueue chan url.URL
type ResponseQueue chan *Message

// AsyncHttpFetcher implements Fetcher
type AsyncHttpFetcher struct {
	//AsyncHttpFetcher is an Asynchronous Worker
	*AsyncWorker

	requestQueue  *RequestQueue
	responseQueue *ResponseQueue

	client HTPPClient
}

func NewAsyncHttpFetcher() *AsyncHttpFetcher {
	reqQueue := make(RequestQueue)
	resQueue := make(ResponseQueue)
	a := &AsyncHttpFetcher{
		AsyncWorker: &AsyncWorker{
			Name: "Fetcher",
		},

		client:        &http.Client{},
		requestQueue:  &reqQueue,
		responseQueue: &resQueue,
	}
	a.AsyncWorker.RunFunc = a.Run

	return a
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

func (a *AsyncHttpFetcher) Worker() Worker {
	return a.AsyncWorker
}

func (a *AsyncHttpFetcher) Run() error {
	a.AsyncWorker.SetState(RUNNING)
	for {
		a.AsyncWorker.SetState(WAITING)
		select {

		// A request is received
		case req := <-*a.requestQueue:
			a.AsyncWorker.SetState(RUNNING)
			res, _ := a.get(req)
			*a.responseQueue <- &Message{
				Request:  &req,
				Response: res,
			}

		// A quit has been received, Stop has been invoked
		case <-a.AsyncWorker.quit:
			a.AsyncWorker.SetState(STOPPED)
			return nil
		default:
		}
	}
}
