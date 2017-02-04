package crawl

import (
	"fmt"
	"net/url"
	"time"

	"github.com/antoniou/go-crawler/fetch"
	"github.com/antoniou/go-crawler/sitemap"
	"github.com/antoniou/go-crawler/track"
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
func NewAsyncHTTPCrawler(f fetch.Fetcher, t track.Tracker, workers []fetch.Worker) *AsyncHTTPCrawler {
	c := &AsyncHTTPCrawler{
		fetcher: f,
		tracker: t,
		workers: append(workers, f.Worker(), t.Worker()),
	}

	return c
}

// AsyncHTTPCrawler is an implementation of the
// Crawler interface. It contains a fetcher that
// initiates the crawling and zero or more workers
// that perform the processing
type AsyncHTTPCrawler struct {
	fetcher fetch.Fetcher
	tracker track.Tracker
	workers []fetch.Worker
}

// Crawl is the main entrypoint to crawling a domain (url).
// Crawl returns a Sitemapper that can later be used to create a
// represenation of the crawled site.
// It returns an error in case the crawl url is invalid
func (c *AsyncHTTPCrawler) Crawl(url *url.URL) (sitemap.Sitemapper, error) {
	stmp := sitemap.New()
	c.tracker.SetSitemapper(stmp)

	for _, worker := range c.workers {
		fmt.Printf("Starting worker of type %v\n", worker.Type())
		go worker.Run()
	}

	err := c.fetcher.Fetch(url)
	if err != nil {
		return nil, err
	}
	c.join()

	return stmp, nil
}

// Wait for all workers to be in state WAITING. This
// will indicate that work is done
func (c *AsyncHTTPCrawler) join() {
	for {
		time.Sleep(1 * time.Second)
		state := fetch.WAITING

		for _, worker := range c.workers {
			state += worker.State()
		}

		if state == fetch.WAITING {
			return
		}

	}
}
