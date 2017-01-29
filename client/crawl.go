package client

import (
	"fmt"
	"log"

	"github.com/antoniou/go-crawler/crawl"
	"github.com/antoniou/go-crawler/fetch"
	"github.com/antoniou/go-crawler/parse"
)

// CrawlCommand is a Command implementation that
// crawls
type CrawlCommand struct {
	BaseCommand
}

// NewCrawlCommand returns a CrawlCommand instance
func NewCrawlCommand() *CrawlCommand {

	return &CrawlCommand{
		BaseCommand{
			Name:        "crawl",
			Usage:       "Crawl a URL",
			Description: "Crawl a URL.",
			ArgsUsage:   "<URL>",
		},
	}
}

func (c *CrawlCommand) Run(args []string) error {
	if len(args) == 0 {
		log.Fatalf("The %s command expects at least one argument", c.Name)
	}
	fetcher := fetch.NewAsyncHttpFetcher()
	parser := parse.NewAsynchHttpParser()
	crawler := crawl.New(fetcher, parser)
	crawler.Crawl(args[0])
	for {
		res := <-*parser.ResponseQueue
		fmt.Printf("%s\n", res)

	}
}
