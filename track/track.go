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
	Run(parser parse.Parser, fetcher fetch.Fetcher) error
}

type AsynchHttpTracker struct {
	filter     *bloom.BloomFilter
	sitemapper sitemap.Sitemapper
}

func New() *AsynchHttpTracker {
	filter := bloom.New(20000, 5)
	t := &AsynchHttpTracker{
		filter:     filter,
		sitemapper: sitemap.New(),
	}
	return t
}

func (t *AsynchHttpTracker) Run(parser parse.Parser, fetcher fetch.Fetcher) error {
	fmt.Println("Tracker: Starting...")
	for {
		fmt.Println("Tracker: Waiting for input from Parser")
		res, _ := parser.Retrieve()
		if !t.filter.TestAndAddString(res) {

			// Adding to sitemapper
			t.sitemapper.Add(res)
			url, err := url.ParseRequestURI(res)
			if err != nil {
				return err
			}
			fmt.Printf("Tracker: Requesting to fetch %s from Fetcher\n", url)
			go fetcher.Fetch(url)
		} else {
			fmt.Printf("Tracker: Url %s is already in bloom filter\n", res)
		}
	}

}
