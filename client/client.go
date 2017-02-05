package client

import (
	"fmt"
	"os"

	"github.com/antoniou/go-crawler/crawl"
	"github.com/antoniou/go-crawler/sitemap"
	"github.com/antoniou/go-crawler/util"
	"github.com/urfave/cli"
)

//Client represents a command line client
type Client struct {
	app *cli.App
}

// Run starts the command line client
func (client *Client) Run(arguments []string) error {
	return client.app.Run(arguments)
}

// New is a Client constructor
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
		// Logger with verbose logging enabled/disabled
		_ = util.Logger(c.Bool("verbose"))
		if len(c.Args()) == 0 {
			client.app.Commands[0].Run(c)
			return fmt.Errorf("Expects at least one argument")
		}
		return client.crawl(c)
	}

	client.app = app
	return client
}

// crawl initiates the crawling steps
// Requires one or more valid url strings
func (client *Client) crawl(c *cli.Context) error {
	args := c.Args()

	seedURL, err := util.NormalizeStringURL(args[0])
	if err != nil {
		return err
	}

	crawler := crawl.NewAsyncHTTPCrawler(seedURL)
	stmp, err := crawler.Crawl()
	if err != nil {
		return err
	}

	outfile := c.String("o")
	return client.export(outfile, stmp)
}

// export sitemap stmp to new file outfile
func (client *Client) export(outfile string, stmp sitemap.Sitemapper) error {
	f, err := os.Create(outfile)
	if err != nil {
		return err
	}

	err = sitemap.NewExporter(f).Export(stmp)
	if err != nil {
		return err
	}
	fmt.Printf("Sitemap exported to %s\n", outfile)
	return nil
}
