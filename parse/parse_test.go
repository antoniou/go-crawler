package parse

import (
	"net/url"
	"testing"

	"github.com/antoniou/go-crawler/fetch"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock the fetcher that is going to be passed to the Parser tests
type mockFetcher struct {
	mock.Mock
	fetch.Fetcher

	responseQueue *fetch.ResponseQueue
}

func (f *mockFetcher) Run(url *url.URL) error {
	return nil
}

func (f *mockFetcher) ResponseChannel() (Response *fetch.ResponseQueue) {
	return f.responseQueue
}

func (f *mockFetcher) Worker() fetch.Worker {
	return &fetch.AsyncWorker{}
}

func (f *mockFetcher) Fetch(url *url.URL) error {
	// f := os.Open(filepath.Join("testdata", url.String()))
	// f.Read(b)
	return nil
}

type ParseTestSuite struct {
	suite.Suite
	seedURL *url.URL
}

func NewTestParser(seedURL *url.URL, fetcher fetch.Fetcher) *AsyncHTTPParser {
	resQueue := make(ResponseQueue)
	a := &AsyncHTTPParser{
		AsyncWorker: fetch.NewAsyncWorker("Parser"),

		fetcher:       fetcher,
		ResponseQueue: &resQueue,
		seed:          seedURL,
	}
	a.AsyncWorker.RunFunc = a.Run
	go a.Worker().Run()
	return a
}

func NewMockFetcher() fetch.Fetcher {
	fetcherResQueue := make(fetch.ResponseQueue)
	f := &mockFetcher{
		responseQueue: &fetcherResQueue,
	}
	return f
}

func (suite *ParseTestSuite) SetupTest() {
	suite.seedURL, _ = url.ParseRequestURI("http://example.com")
}

func (suite *ParseTestSuite) TestFetchValidAndInvalidResponse() {
	f := NewMockFetcher()
	p := NewTestParser(suite.seedURL, f)
	f.Fetch(p.seed)

}

func (suite *ParseTestSuite) TestStopParser() {
	f := NewMockFetcher()
	p := NewTestParser(suite.seedURL, f)
	assert.Equal(suite.T(), fetch.WAITING, p.Worker().State())

	p.Stop()
	assert.Equal(suite.T(), fetch.STOPPED, p.Worker().State())
}

func (suite *ParseTestSuite) TestNewAsyncHTTPParserConstructor() {
	f := NewMockFetcher()
	p := NewAsyncHTTPParser(suite.seedURL, f)
	assert.Implements(suite.T(), (*Parser)(nil), p)
	assert.NotNil(suite.T(), p)
}

func TestParseTestSuite(t *testing.T) {
	suite.Run(t, new(ParseTestSuite))
}
