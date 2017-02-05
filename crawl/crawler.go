package crawl

import (
	"fmt"
	"net/url"
	"time"

	"github.com/antoniou/go-crawler/sitemap"
	"github.com/antoniou/go-crawler/util"
)

// A Crawler crawls a domain and returns
// a representation of the crawled domain
type Crawler interface {
	//Crawl is the main entrypoint to
	//crawling a domain (url)
	Crawl(url string) (sitemap.Sitemapper, error)
}

// NewAsyncHTTPCrawler is a constructor. It takes in a Fetcher
// that will start the crawl and zero or more workers that will
// process the response and create a Sitemap
func NewAsyncHTTPCrawler(seedURL *url.URL) *AsyncHTTPCrawler {

	fetcher := NewAsyncHTTPFetcher()
	parser := NewAsyncHTTPParser(seedURL, fetcher)
	tracker := NewAsyncHttpTracker(fetcher, parser)
	return &AsyncHTTPCrawler{
		seedURL: seedURL,
		fetcher: fetcher,
		tracker: tracker,
		workers: []Worker{
			parser.Worker(),
			fetcher.Worker(),
			tracker.Worker(),
		},
	}
}

// AsyncHTTPCrawler is an implementation of the
// Crawler interface. It contains a fetcher that
// initiates the crawling and zero or more workers
// that perform the processing
type AsyncHTTPCrawler struct {
	fetcher Fetcher
	tracker Tracker
	workers []Worker
	seedURL *url.URL
}

// Crawl is the main entrypoint to crawling a domain (url).
// Crawl returns a Sitemapper that can later be used to create a
// represenation of the crawled site.
// It returns an error in case the crawl url is invalid
func (c *AsyncHTTPCrawler) Crawl() (sitemap.Sitemapper, error) {
	// Create an empty sitemap
	stmp := sitemap.NewGraphSitemap()
	// Pass it to the tracker
	c.tracker.SetSitemapper(stmp)

	for _, worker := range c.workers {
		util.Printf("Starting worker of type %v\n", worker.Type())
		go worker.Run()
	}

	fmt.Printf("Starting crawling of %v\n", c.seedURL)
	err := c.fetcher.Fetch(c.seedURL)
	if err != nil {
		return nil, err
	}

	return stmp, c.join()
}

// Wait for all workers to be in state WAITING. This
// will indicate that work is done
func (c *AsyncHTTPCrawler) join() error {
	for {
		time.Sleep(500 * time.Millisecond)
		state := WAITING

		for _, worker := range c.workers {
			state += worker.State()
		}

		if state == WAITING {
			return nil
		}

	}
}
