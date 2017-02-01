package fetch

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	WAITING = 0
	STOPPED = 1
	RUNNING = 2
)

// Fetcher is an Asynchronous interface
type Fetcher interface {
	// Fetch provides work to the Fetcher, in the
	// form of a URL to process
	Fetch(url *url.URL) error

	// Retrieve provides results back from the Fetcher
	// in the form of a Response
	Retrieve() (Response *Message, err error)

	// Retrieve Worker
	Worker() Worker
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

type Message struct {
	Request  *url.URL
	Response io.ReadCloser
}

type RequestQueue chan url.URL
type ResponseQueue chan *Message

// AsyncWorker is an interface for
type Worker interface {
	// Run starts the Asynchronous worker
	Run() error

	// State returns the state the worker is in:
	// RUNNING - processing work
	// WAITING - Waits for work
	// STOPPED - Not running
	State() uint8
	SetState(state uint8)

	// Returns worker name
	// Possible names are:
	// Fetcher
	// Parser
	// Tracker
	// Sitemapper
	Type() string
}

type AsyncWorker struct {
	RunFunc func() error

	state uint8
	quit  chan uint8
	Name  string
}

func (w *AsyncWorker) Run() error {
	return w.RunFunc()
}

func (w *AsyncWorker) Stop() error {
	w.quit <- 0
	return nil
}

func (w *AsyncWorker) State() uint8 {
	return w.state
}

func (w *AsyncWorker) SetState(state uint8) {
	w.state = state
}

func (w *AsyncWorker) Type() string {
	return w.Name
}

// AsyncHttpFetcher implements Fetcher
type AsyncHttpFetcher struct {
	//AsyncHttpFetcher is an Asynchronous Worker
	*AsyncWorker

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
			fmt.Printf("Fetcher: Going to fetch url %s\n", req.String())
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
