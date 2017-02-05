package crawl

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/goware/urlx"
	"golang.org/x/net/html"
)

type ParserResponseQueue chan *ParseMessage

// Parser is an Asynchronous interface
type Parser interface {
	// Run starts the Parser
	Run() error

	// Retrieve provides results back from the Parser
	// in the form of Urls
	Retrieve() (m *ParseMessage, err error)

	// Retrieve Worker
	Worker() Worker
}

type AsyncHTTPParser struct {
	//Fetcher is an Asynchronous Worker
	*AsyncWorker

	fetcher             Fetcher
	ParserResponseQueue *ParserResponseQueue
	seed                *url.URL
}

type ParseMessage struct {
	Request  *url.URL
	Response *string
}

func NewAsyncHTTPParser(seedURL *url.URL, fetcher Fetcher) *AsyncHTTPParser {
	resQueue := make(ParserResponseQueue)
	a := &AsyncHTTPParser{
		AsyncWorker: NewAsyncWorker("Parser"),

		fetcher:             fetcher,
		ParserResponseQueue: &resQueue,
		seed:                seedURL,
	}
	a.AsyncWorker.RunFunc = a.Run
	return a
}

func (p *AsyncHTTPParser) Run() error {
	p.AsyncWorker.SetState(RUNNING)
	for {
		p.AsyncWorker.SetState(WAITING)
		select {
		case res := <-*p.fetcher.ResponseChannel():
			if err := p.handleResponse(res); err != nil {
				continue
			}
		case <-p.AsyncWorker.Quit:
			p.Worker().SetState(STOPPED)
			return nil
		}

	}
}

func (p *AsyncHTTPParser) handleResponse(res *FetchMessage) error {
	if res.Error != nil {
		if res.Request.String() == p.seed.String() {
			p.Stop()
		}
		log.Printf("Could not get %s: %v", res.Request.String(), res.Error)
		return res.Error
	}
	p.AsyncWorker.SetState(RUNNING)
	p.extractLinks(res)
	return nil
}

func (p *AsyncHTTPParser) extractLinks(res *FetchMessage) error {
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
			isRelative := strings.HasPrefix(url, "/")
			if isRelative {
				url = fmt.Sprintf("%s://%s%s", p.seed.Scheme, p.seed.Host, url)
			}
			normURL := p.normalise(url)
			hasProto := strings.Index(normURL, "http") == 0
			inSeedDomain := strings.Index(normURL, p.seed.String()) == 0
			if hasProto && inSeedDomain {
				*p.ParserResponseQueue <- &ParseMessage{
					Request:  res.Request,
					Response: &normURL,
				}
			}
		}
	}

	return nil
}

func (p *AsyncHTTPParser) Retrieve() (m *ParseMessage, err error) {
	m = <-*p.ParserResponseQueue
	return m, nil

}

func (a *AsyncHTTPParser) Worker() Worker {
	return a.AsyncWorker
}

// Helper function to bring Url to its normalised form
// by removing querystrings and reconstructing absolute
// path URLs
func (p *AsyncHTTPParser) normalise(path string) string {
	normURL := path
	parsedURL, err := url.ParseRequestURI(normURL)
	if err != nil {
		return ""
	}

	normURL = strings.TrimSuffix(path, "?"+parsedURL.RawQuery)
	normalized, _ := urlx.NormalizeString(normURL)
	return normalized
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
