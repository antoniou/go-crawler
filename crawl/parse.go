package crawl

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"github.com/antoniou/go-crawler/util"
	"golang.org/x/net/html"
)

type parserResponseQueue chan *ParseMessage

// Parser is an Asynchronous interface
type Parser interface {
	// ResponseChannel is a Getter returning
	// the Parser's Channel  that consumers
	// should be receiving results from
	ResponseChannel() (responseQueue *parserResponseQueue)

	// Retrieve Worker
	Worker() Worker
}

type AsyncHTTPParser struct {
	//Fetcher is an Asynchronous Worker
	*AsyncWorker

	fetcher             Fetcher
	parserResponseQueue *parserResponseQueue
	seed                *url.URL
}

type ParseMessage struct {
	Request  *url.URL
	Response *url.URL
}

func NewAsyncHTTPParser(seedURL *url.URL, fetcher Fetcher) *AsyncHTTPParser {
	resQueue := make(parserResponseQueue)
	a := &AsyncHTTPParser{
		AsyncWorker: NewAsyncWorker("Parser"),

		fetcher:             fetcher,
		parserResponseQueue: &resQueue,
		seed:                seedURL,
	}
	a.AsyncWorker.RunFunc = a.Run
	return a
}

// Run starts a loop that waits for requests
// or the quit signal. Run will be interrupted
// once the Stop method is used
func (p *AsyncHTTPParser) Run() error {
	p.AsyncWorker.SetState(RUNNING)
	for {
		p.AsyncWorker.SetState(WAITING)
		select {
		case res := <-*p.fetcher.ResponseChannel():
			if err := p.handleResponse(res); err != nil {
				if p.seed.String() == res.Request.String() {
					return err
				}
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
			normURL, err := util.NormalizeStringURL(url)
			if err != nil {
				util.Printf("Parser: Error while normalizing %v: %v", url, err)
				continue
			}
			hasProto := strings.Index(normURL.Scheme, "http") == 0
			inSeedDomain := strings.Index(normURL.String(), p.seed.String()) == 0
			if hasProto && inSeedDomain {
				util.Printf("Parser: Passing url %v to Tracker", normURL)
				*p.parserResponseQueue <- &ParseMessage{
					Request:  res.Request,
					Response: normURL,
				}
			}
		}
	}

	return nil
}

func (p *AsyncHTTPParser) ResponseChannel() *parserResponseQueue {
	return p.parserResponseQueue
}

// Worker Returns the embedded AsyncWorker struct
// which is used to Run and Stop the Parser worker
func (p *AsyncHTTPParser) Worker() Worker {
	return p.AsyncWorker
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
