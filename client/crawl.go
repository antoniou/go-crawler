package client

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/antoniou/go-crawler/crawl"
	"github.com/antoniou/go-crawler/fetch"
	"github.com/antoniou/go-crawler/parse"
	"github.com/antoniou/go-crawler/sitemap"
	"github.com/antoniou/go-crawler/track"
	"github.com/goware/urlx"
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

// Run executes the Crawl command.
// Requires one or more valid url strings
func (c *CrawlCommand) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("The %s command expects at least one argument", c.Name)
	}

	normURL, _ := urlx.NormalizeString(args[0])
	url, err := url.ParseRequestURI(normURL)

	if err != nil {
		return err
	}

	fetcher := fetch.NewAsyncHTTPFetcher()
	parser := parse.NewAsyncHTTPParser(url, fetcher)
	tracker := track.New(fetcher, parser)

	crawler := crawl.NewAsyncHTTPCrawler(
		fetcher,
		tracker,
		[]fetch.Worker{
			parser.Worker(),
		},
	)

	stmp, err := crawler.Crawl(url)
	if err != nil {
		log.Println(err)
		return err
	}

	f, err := os.Create("result.txt")
	if err != nil {
		return err
	}

	err = sitemap.NewExporter(f).Export(stmp)
	return err
}
