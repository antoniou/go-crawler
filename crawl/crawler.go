package crawl

import (
	"fmt"
	"net/url"
	"time"

	"github.com/antoniou/go-crawler/fetch"
	"github.com/antoniou/go-crawler/parse"
	"github.com/antoniou/go-crawler/track"
)

type Crawler interface {
	//Crawl is the main entrypoint to
	//crawling a domain
	Crawl(url string) error
}

func New(f fetch.Fetcher, workers []fetch.Worker) *AsyncHttpCrawler {
	c := &AsyncHttpCrawler{
		fetcher: f,
		workers: workers,
	}

	return c
}

type AsyncHttpCrawler struct {
	fetcher fetch.Fetcher
	parser  parse.Parser
	tracker track.Tracker
	workers []fetch.Worker
}

func (c *AsyncHttpCrawler) Crawl(url *url.URL) error {
	for _, worker := range c.workers {
		fmt.Printf("Starting worker of type %v\n", worker.Type())
		go worker.Run()
	}
	err := c.fetcher.Fetch(url)
	if err != nil {
		return err
	}
	c.join()

	return nil
}

func (c *AsyncHttpCrawler) join() {
	for {
		time.Sleep(1 * time.Second)
		var state uint8
		state = fetch.WAITING
		for _, worker := range c.workers {
			state += worker.State()
		}
		if state == fetch.WAITING {
			fmt.Println("All workers are in waiting state, crawling complete")
			return
		}

	}
}
