package client

import (
	"log"
	"net/url"

	"github.com/antoniou/go-crawler/crawl"
	"github.com/antoniou/go-crawler/fetch"
	"github.com/antoniou/go-crawler/parse"
	"github.com/antoniou/go-crawler/track"
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

	url, err := url.ParseRequestURI(args[0])
	if err != nil {
		log.Println(err)
		return err
	}

	fetcher := fetch.NewAsyncHTTPFetcher()
	parser := parse.NewAsyncHTTPParser(url, fetcher)
	tracker := track.New(fetcher, parser)

	crawler := crawl.NewAsyncHTTPCrawler(
		fetcher,
		[]fetch.Worker{
			parser.Worker(),
			tracker.Worker(),
		},
	)

	_, err = crawler.Crawl(url)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
