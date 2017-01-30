package parse

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/antoniou/go-crawler/fetch"
	"golang.org/x/net/html"
)

type ResponseQueue chan string

// Parser is an Asynchronous interface
type Parser interface {
	// Run starts the Parser with a specific
	// fetcher that will supply http Responses to parse
	Run(fetcher fetch.Fetcher) error

	// Retrieve provides results back from the Parser
	// in the form of Urls
	Retrieve() (url string, err error)
}

type AsynchHttpParser struct {
	ResponseQueue *ResponseQueue
	seed          *url.URL
}

func NewAsynchHttpParser(seedUrl *url.URL) *AsynchHttpParser {
	resQueue := make(ResponseQueue)
	a := &AsynchHttpParser{
		ResponseQueue: &resQueue,
		seed:          seedUrl,
	}
	return a
}

func (p *AsynchHttpParser) Run(fetcher fetch.Fetcher) error {

	for {
		fmt.Println("Parser: Waiting for response from fetcher...")
		res, _ := fetcher.Retrieve()

		p.extractLinks(res)
	}
}

func (p *AsynchHttpParser) extractLinks(res *fetch.Message) error {
	z := html.NewTokenizer(res.Response)
	done := false
	for {
		if done {
			break
		}
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
			done = true
			break
		case tt == html.StartTagToken:
			t := z.Token()

			// Check if the token is an <a> tag
			isAnchor := t.Data == "a"
			if !isAnchor {
				continue
			}

			// Extract the href value, if there is one
			ok, url := getHref(t)
			if !ok {
				continue
			}

			// Make sure the url begines in http**
			hasProto := strings.Index(url, "http") == 0
			inSeedDomain := strings.Index(url, p.seed.String()) == 0
			isRelative := strings.HasPrefix(url, "/")
			if isRelative {
				url = fmt.Sprintf("%s%s", p.seed.String(), url)
			}

			if hasProto && inSeedDomain || isRelative {
				*p.ResponseQueue <- url
			}
		}
	}

	return nil
}

func (p *AsynchHttpParser) Retrieve() (url string, err error) {
	url = <-*p.ResponseQueue
	fmt.Printf("Parser: Passing %s to Tracker\n", url)
	return url, nil
}

// Helper function to pull the href attribute from a Token
func getHref(t html.Token) (ok bool, href string) {
	// Iterate over all of the Token's attributes until we find an "href"
	for _, a := range t.Attr {
		if a.Key == "href" {
			href = a.Val
			ok = true
		}
	}
	// "bare" return will return the variables (ok, href) as defined in
	// the function definition
	return
}
