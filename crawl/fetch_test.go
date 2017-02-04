package crawl

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// Define a mock struct to be used in your unit tests of myFunc.
type mockHTTPClient struct {
	mock.Mock
}

func (m *mockHTTPClient) Get(url string) (resp *http.Response, err error) {

	if strings.Contains(url, "nonexistingwebsite") {
		return nil, fmt.Errorf("no such host")
	} else if strings.Contains(url, "404") {
		return &http.Response{
			Status: "404",
		}, nil
	}

	response := &http.Response{
		Status: "200",
	}
	return response, nil
}

type FetchTestSuite struct {
	suite.Suite
}

func NewTestFetcher() *AsyncHTTPFetcher {
	reqQueue := make(RequestQueue)
	resQueue := make(FetchResponseQueue)
	a := &AsyncHTTPFetcher{
		AsyncWorker: NewAsyncWorker("MockFetcher"),

		client:        &mockHTTPClient{},
		requestQueue:  &reqQueue,
		responseQueue: &resQueue,
	}
	a.AsyncWorker.RunFunc = a.Run
	go a.AsyncWorker.Run()
	return a
}

func (suite *FetchTestSuite) TestFetchValidAndInvalidResponse() {
	f := NewTestFetcher()

	// Valid Request
	uri, _ := url.ParseRequestURI("https://validurl.com")
	f.Fetch(uri)
	m := <-*f.ResponseChannel()

	assert.Equal(suite.T(), m.Request.String(), uri.String())
	assert.NoError(suite.T(), m.Error)
	assert.Equal(suite.T(), "200", m.Response.Status)

	// Invalid Scheme
	uri, _ = url.ParseRequestURI("ftp://invalidurl")
	err := f.Fetch(uri)
	assert.Error(suite.T(), err)

	// Cannot resolve website
	uri, _ = url.ParseRequestURI("http://nonexistingwebsite.com")
	f.Fetch(uri)
	m = <-*f.ResponseChannel()
	assert.Error(suite.T(), m.Error)
	assert.Nil(suite.T(), m.Response)

	// Website returns 404
	uri, _ = url.ParseRequestURI("http://iama404response.com")
	f.Fetch(uri)
	m = <-*f.ResponseChannel()
	assert.Equal(suite.T(), "404", m.Response.Status)
	assert.Nil(suite.T(), m.Error)

}

func (suite *FetchTestSuite) TestStopFetcher() {
	f := NewTestFetcher()
	assert.Equal(suite.T(), WAITING, f.Worker().State())
	f.Stop()
	assert.Equal(suite.T(), STOPPED, f.Worker().State())

	uri, _ := url.ParseRequestURI("https://validurl.com")
	err := f.Fetch(uri)
	assert.Error(suite.T(), err)
	assert.Contains(suite.T(), err.Error(), "is in state stopped")
}

func (suite *FetchTestSuite) TestNewAsyncHTTPFetcherConstructor() {
	f := NewAsyncHTTPFetcher()
	assert.Implements(suite.T(), (*Fetcher)(nil), f)
	assert.NotNil(suite.T(), f)
}

func TestFetchTestSuite(t *testing.T) {
	suite.Run(t, new(FetchTestSuite))
}
