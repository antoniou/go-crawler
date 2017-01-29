package crawl

import (
	"fmt"
	"net/url"

	"github.com/antoniou/go-crawler/fetch"
	"github.com/antoniou/go-crawler/parse"
	"github.com/antoniou/go-crawler/track"
)

type Crawler interface {
	//Crawl is the main entrypoint to
	//crawling a domain
	Crawl(url string) error
}

func New(f fetch.Fetcher, p parse.Parser, t track.Tracker) *AsyncHttpCrawler {
	c := &AsyncHttpCrawler{
		fetcher: f,
		parser:  p,
		tracker: t,
	}
	return c
}

type AsyncHttpCrawler struct {
	fetcher fetch.Fetcher
	parser  parse.Parser
	tracker track.Tracker
}

func (c *AsyncHttpCrawler) Crawl(url *url.URL) error {
	go c.parser.Run(c.fetcher)
	go c.tracker.Run(c.parser, c.fetcher)
	go c.fetcher.Run()
	fmt.Printf("Crawler: Start crawling with url %v\n", url)
	go c.fetcher.Fetch(url)

	return nil
}
