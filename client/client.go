package client

import (
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/antoniou/go-crawler/crawl"
	"github.com/antoniou/go-crawler/sitemap"
	"github.com/antoniou/go-crawler/util"
	"github.com/goware/urlx"
	"github.com/urfave/cli"
)

//Client represents a command line client
type Client struct {
	app *cli.App
}

func (client *Client) Run(arguments []string) error {
	return client.app.Run(arguments)
}

func New() (client *Client) {
	client = new(Client)
	app := cli.NewApp()
	app.Name = "go-crawler"
	app.Usage = "Crawl a site"
	app.UsageText = "crawl [options] url"
	app.Version = "0.1.0"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "o",
			Value: "result.out",
			Usage: "Output file",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Verbose mode",
		},
	}

	app.Action = func(c *cli.Context) error {
		return client.Crawl(c)
	}

	client.app = app
	return client
}

// Crawl initiates the crawling steps
// Requires one or more valid url strings
func (client *Client) Crawl(c *cli.Context) error {
	args := c.Args()

	if len(args) == 0 {
		return fmt.Errorf("Expects at least one argument")
	}

	normURL, _ := urlx.NormalizeString(args[0])
	url, err := url.ParseRequestURI(normURL)

	if err != nil {
		return err
	}

	_ = util.Logger(c.Bool("verbose"))
	fetcher := crawl.NewAsyncHTTPFetcher()
	parser := crawl.NewAsyncHTTPParser(url, fetcher)
	tracker := crawl.New(fetcher, parser)

	crawler := crawl.NewAsyncHTTPCrawler(
		fetcher,
		tracker,
		[]crawl.Worker{
			parser.Worker(),
		},
	)

	stmp, err := crawler.Crawl(url)
	if err != nil {
		log.Println(err)
		return err
	}

	outfile := c.String("o")
	f, err := os.Create(outfile)
	if err != nil {
		return err
	}

	err = sitemap.NewExporter(f).Export(stmp)
	if err != nil {
		return err
	}

	fmt.Printf("Sitemap exported to %s\n", outfile)
	return err
}
