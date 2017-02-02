package parse

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/antoniou/go-crawler/fetch"
	"golang.org/x/net/html"
)

type ResponseQueue chan *Message

// Parser is an Asynchronous interface
type Parser interface {
	// Run starts the Parser
	Run() error

	// Retrieve provides results back from the Parser
	// in the form of Urls
	Retrieve() (m *Message, err error)

	// Retrieve Worker
	Worker() fetch.Worker
}

type AsyncHttpParser struct {
	//Fetcher is an Asynchronous Worker
	*fetch.AsyncWorker

	fetcher       fetch.Fetcher
	ResponseQueue *ResponseQueue
	seed          *url.URL
}

type Message struct {
	Request  *url.URL
	Response *string
}

func NewAsyncHttpParser(seedUrl *url.URL, fetcher fetch.Fetcher) *AsyncHttpParser {
	resQueue := make(ResponseQueue)
	a := &AsyncHttpParser{
		AsyncWorker: &fetch.AsyncWorker{
			Name: "Parser",
		},

		fetcher:       fetcher,
		ResponseQueue: &resQueue,
		seed:          seedUrl,
	}
	a.AsyncWorker.RunFunc = a.Run
	return a
}

func (p *AsyncHttpParser) Run() error {
	p.AsyncWorker.SetState(fetch.RUNNING)
	for {
		p.AsyncWorker.SetState(fetch.WAITING)
		res, err := p.fetcher.Retrieve()
		if err != nil {
			log.Printf("Could not get %s: %v", res.Request.String(), err)
			continue
		}
		p.AsyncWorker.SetState(fetch.RUNNING)
		p.extractLinks(res)
	}
}

func (p *AsyncHttpParser) extractLinks(res *fetch.Message) error {
	z := html.NewTokenizer(res.Response.Body)
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
			if (hasProto && inSeedDomain) || isRelative {
				normURL := p.normalise(url)
				*p.ResponseQueue <- &Message{
					Request:  res.Request,
					Response: &normURL,
				}
			}
		}
	}

	return nil
}

func (p *AsyncHttpParser) Retrieve() (m *Message, err error) {
	m = <-*p.ResponseQueue
	fmt.Printf("Parser: %s -> %s\n", m.Request.String(), *m.Response)
	return m, nil

}

func (a *AsyncHttpParser) Worker() fetch.Worker {
	return a.AsyncWorker
}

// Helper function to bring Url to its normalised form
// by removing querystrings and reconstructing absolute
// path URLs
func (p *AsyncHttpParser) normalise(path string) string {
	normURL := path
	parsedURL, err := url.ParseRequestURI(normURL)
	if err != nil {
		log.Fatal(err)
	}

	return strings.TrimSuffix(path, "?"+parsedURL.RawQuery)
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
