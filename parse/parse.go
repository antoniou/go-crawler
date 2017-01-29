package parse

import (
	"fmt"
	"strings"

	"github.com/antoniou/go-crawler/fetch"
	"golang.org/x/net/html"
)

type ResponseQueue chan string

// Parser is an Asynchronous interface
type Parser interface {
	Run(fetcher fetch.Fetcher) error
}

type AsynchHttpParser struct {
	ResponseQueue *ResponseQueue
}

func NewAsynchHttpParser() *AsynchHttpParser {
	resQueue := make(ResponseQueue)
	a := &AsynchHttpParser{
		ResponseQueue: &resQueue,
	}
	return a
}

func (p *AsynchHttpParser) Run(fetcher fetch.Fetcher) error {

	res, _ := fetcher.Retrieve()
	fmt.Println("Parser: Got result!")
	z := html.NewTokenizer(res)

	for {
		tt := z.Next()

		switch {
		case tt == html.ErrorToken:
			// End of the document, we're done
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
			if hasProto {
				*p.ResponseQueue <- url
			}
		}
	}
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
