package crawl

import (
	"fmt"

	"github.com/antoniou/go-crawler/fetch"
	"github.com/antoniou/go-crawler/parse"
)

type Crawler interface {
	//Crawl is the main entrypoint to
	//crawling a domain
	Crawl(url string) error
}

func New(f fetch.Fetcher, p parse.Parser) *AsyncHttpCrawler {
	c := &AsyncHttpCrawler{
		fetcher: f,
		parser:  p,
	}
	return c
}

type AsyncHttpCrawler struct {
	fetcher fetch.Fetcher
	parser  parse.Parser
}

func (c *AsyncHttpCrawler) Crawl(url string) error {
	go c.parser.Run(c.fetcher)
	fmt.Printf("Crawler: Start crawling with url %s\n", url)
	c.fetcher.Fetch(url)

	return nil
}
