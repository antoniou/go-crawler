package track

import (
	"fmt"
	"net/url"

	"github.com/antoniou/go-crawler/fetch"
	"github.com/antoniou/go-crawler/parse"
	"github.com/antoniou/go-crawler/sitemap"
	"github.com/willf/bloom"
)

type Tracker interface {
	// Retrieve Worker
	Worker() fetch.Worker

	// SetSitemapper provides the Tracker with
	// a Sitemapper. The Tracker is responsible for
	// building the providing the Sitemapper with
	// new URL data
	SetSitemapper(sitemap.Sitemapper)
}

type AsyncHttpTracker struct {
	//Tracker is an Asynchronous Worker
	*fetch.AsyncWorker

	filter     *bloom.BloomFilter
	fetcher    fetch.Fetcher
	parser     parse.Parser
	sitemapper sitemap.Sitemapper
}

func New(fetcher fetch.Fetcher, parser parse.Parser) *AsyncHttpTracker {
	filter := bloom.New(20000, 5)
	t := &AsyncHttpTracker{
		AsyncWorker: &fetch.AsyncWorker{
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
	t.Worker().SetState(fetch.RUNNING)
	for {
		t.Worker().SetState(fetch.WAITING)
		res, _ := t.parser.Retrieve()
		t.Worker().SetState(fetch.RUNNING)
		t.handle(res)
	}
}

func (t *AsyncHttpTracker) handle(m *parse.Message) error {
	if t.filter.TestAndAddString(*m.Response) {
		return nil
	}

	fmt.Printf("Tracker: Adding %s to sitemap\n", *m.Response)
	t.sitemapper.Add(m.Request.String(), *m.Response)
	url, err := url.ParseRequestURI(*m.Response)
	if err != nil {
		return err
	}

	go t.fetcher.Fetch(url)
	return nil
}

// SetSitemapper provides the Tracker with
// a Sitemapper. The Tracker is responsible for
// building the providing the Sitemapper with
// new URL data
func (t *AsyncHttpTracker) SetSitemapper(s sitemap.Sitemapper) {
	t.sitemapper = s
}

func (t *AsyncHttpTracker) Worker() fetch.Worker {
	return t.AsyncWorker
}
