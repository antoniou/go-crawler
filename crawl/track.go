package crawl

import (
	"github.com/antoniou/go-crawler/sitemap"
	"github.com/antoniou/go-crawler/util"
	"github.com/willf/bloom"
)

// A Tracker is an Asynchronous worker interface
// that is responsible for receiving URLs from the
type Tracker interface {
	// SetSitemapper provides the Tracker with
	// a Sitemapper. The Tracker is responsible for
	// building the providing the Sitemapper with
	// new URL data.
	SetSitemapper(sitemap.Sitemapper)

	// Retrieve Worker
	Worker() Worker
}

// An AsyncHttpTracker is an Asynchronous worker struct
// that is responsible for receiving URLs from a Parser
// and passing the uncrawled URLs to the Fetcher
type AsyncHttpTracker struct {
	//Tracker is an Asynchronous Worker
	*AsyncWorker

	filter     *bloom.BloomFilter
	fetcher    Fetcher
	parser     Parser
	sitemapper sitemap.Sitemapper
}

func NewAsyncHttpTracker(fetcher Fetcher, parser Parser) *AsyncHttpTracker {
	// FIXME Revisit bloom filter size as future work
	filter := bloom.New(20000, 5)
	t := &AsyncHttpTracker{
		AsyncWorker: &AsyncWorker{
			Name: "Tracker",
		},

		filter:  filter,
		fetcher: fetcher,
		parser:  parser,
	}
	t.AsyncWorker.RunFunc = t.Run
	return t
}

func (t *AsyncHttpTracker) Run() error {
	t.Worker().SetState(RUNNING)
	for {
		t.Worker().SetState(WAITING)
		select {
		case res := <-*t.parser.ResponseChannel():
			t.Worker().SetState(RUNNING)
			if err := t.handleResponse(res); err != nil {
				continue
			}
		case <-t.AsyncWorker.Quit:
			t.Worker().SetState(STOPPED)
			return nil
		}
	}
}

func (t *AsyncHttpTracker) handleResponse(m *ParseMessage) error {
	sURL := m.Response.String()
	if t.filter.TestAndAddString(sURL) {
		return nil
	}

	util.Printf("Tracker: Adding %s to sitemap\n", sURL)
	t.sitemapper.Add(m.Request.String(), sURL)

	go t.fetcher.Fetch(m.Response)
	return nil
}

// SetSitemapper provides the Tracker with
// a Sitemapper. The Tracker is responsible for
// building the providing the Sitemapper with
// new URL data
func (t *AsyncHttpTracker) SetSitemapper(s sitemap.Sitemapper) {
	t.sitemapper = s
}

func (t *AsyncHttpTracker) Worker() Worker {
	return t.AsyncWorker
}
