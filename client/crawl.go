package client

import (
	"fmt"
	"log"

	"github.com/antoniou/go-crawler/fetch"
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

	fmt.Print("Starting Async Fetcher...")
	fetcher := fetch.NewAsyncHttpFetcher()
	fmt.Println("Done!")

	fetcher.Fetch(args[0])
	res, _ := fetcher.Retrieve()
	fmt.Printf("Response is %s\n", res)

	return nil
}
