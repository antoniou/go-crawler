package crawl

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Mock the fetcher that is going to be passed to the Parser tests
type mockFetcher struct {
	mock.Mock
	Fetcher

	responseQueue *FetchResponseQueue
}

func (f *mockFetcher) Run(url *url.URL) error {
	return nil
}

func (f *mockFetcher) ResponseChannel() (Response *FetchResponseQueue) {
	return f.responseQueue
}

func (f *mockFetcher) Worker() Worker {
	return &AsyncWorker{}
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

func NewTestParser(seedURL *url.URL, fetcher Fetcher) *AsyncHTTPParser {
	resQueue := make(ParserResponseQueue)
	a := &AsyncHTTPParser{
		AsyncWorker: NewAsyncWorker("Parser"),

		fetcher:             fetcher,
		ParserResponseQueue: &resQueue,
		seed:                seedURL,
	}
	a.AsyncWorker.RunFunc = a.Run
	go a.Worker().Run()
	return a
}

func NewMockFetcher() Fetcher {
	fetcherResQueue := make(FetchResponseQueue)
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
	assert.Equal(suite.T(), WAITING, p.Worker().State())

	p.Stop()
	assert.Equal(suite.T(), STOPPED, p.Worker().State())
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
