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
}

type AsyncHttpTracker struct {
	//Tracker is an Asynchronous Worker
	*fetch.AsyncWorker

	filter     *bloom.BloomFilter
	sitemapper sitemap.Sitemapper
	fetcher    fetch.Fetcher
	parser     parse.Parser
}

func New(fetcher fetch.Fetcher, parser parse.Parser) *AsyncHttpTracker {
	filter := bloom.New(20000, 5)
	t := &AsyncHttpTracker{
		AsyncWorker: &fetch.AsyncWorker{
			Name: "Tracker",
		},

		filter:     filter,
		fetcher:    fetcher,
		parser:     parser,
		sitemapper: sitemap.New(),
	}
	t.AsyncWorker.RunFunc = t.Run
	return t
}

func (t *AsyncHttpTracker) Run() error {
	t.AsyncWorker.SetState(fetch.RUNNING)
	for {
		t.AsyncWorker.SetState(fetch.WAITING)
		res, _ := t.parser.Retrieve()
		t.AsyncWorker.SetState(fetch.RUNNING)
		if !t.filter.TestAndAddString(*res.Response) {

			// Adding to sitemapper
			t.sitemapper.Add(res.Request.String(), *res.Response)
			url, err := url.ParseRequestURI(*res.Response)
			if err != nil {
				return err
			}
			fmt.Printf("Tracker: Requesting to fetch %s from Fetcher\n", url)
			go t.fetcher.Fetch(url)
		} else {
			fmt.Printf("Tracker: Url %s is already in bloom filter\n", *res.Response)
		}
	}
}

func (t *AsyncHttpTracker) Worker() fetch.Worker {
	return t.AsyncWorker
}
