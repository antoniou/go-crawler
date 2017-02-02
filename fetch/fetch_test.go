package fetch

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
type mockHttpClient struct {
	mock.Mock
}

func (m *mockHttpClient) Get(url string) (resp *http.Response, err error) {

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

func NewMockFetcher() *AsyncHTTPFetcher {
	reqQueue := make(RequestQueue)
	resQueue := make(ResponseQueue)
	a := &AsyncHTTPFetcher{
		AsyncWorker: NewAsyncWorker("MockFetcher"),

		client:        &mockHttpClient{},
		requestQueue:  &reqQueue,
		responseQueue: &resQueue,
	}
	a.AsyncWorker.RunFunc = a.Run
	go a.Run()
	return a
}

func (suite *FetchTestSuite) TestFetchValidAndInvalidResponse() {
	f := NewMockFetcher()

	// Valid Request
	uri, _ := url.ParseRequestURI("https://validurl.com")
	f.Fetch(uri)
	m, err := f.Retrieve()

	assert.Equal(suite.T(), m.Request.String(), uri.String())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "200", m.Response.Status)

	// Invalid Scheme
	uri, _ = url.ParseRequestURI("ftp://invalidurl")
	err = f.Fetch(uri)
	assert.Error(suite.T(), err)

	// Cannot resolve website
	uri, _ = url.ParseRequestURI("http://nonexistingwebsite.com")
	f.Fetch(uri)
	m, err = f.Retrieve()
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), m.Response)

	// Website returns 404
	uri, _ = url.ParseRequestURI("http://iama404response.com")
	f.Fetch(uri)
	m, err = f.Retrieve()
	assert.Equal(suite.T(), "404", m.Response.Status)
	assert.Nil(suite.T(), err)

}

func (suite *FetchTestSuite) TestStopFetcher() {
	f := NewMockFetcher()
	assert.Equal(suite.T(), WAITING, f.AsyncWorker.State())
	f.Stop()
	assert.Equal(suite.T(), STOPPED, f.AsyncWorker.State())
}

func TestFetchTestSuite(t *testing.T) {
	suite.Run(t, new(FetchTestSuite))
}
