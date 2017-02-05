package util

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeStringURL(t *testing.T) {
	url, err := NormalizeStringURL("http://example.com")
	assert.Equal(t, url.String(), "http://example.com")
	assert.NoError(t, err)

	url, err = NormalizeStringURL("http://example.com:80")
	assert.Equal(t, url.String(), "http://example.com")
	assert.NoError(t, err)

	url, err = NormalizeStringURL("example.com")
	assert.Equal(t, url.String(), "http://example.com")
	assert.NoError(t, err)

	url, err = NormalizeStringURL("www.example.com")
	assert.Equal(t, url.String(), "http://www.example.com")
	assert.NoError(t, err)

	url, err = NormalizeStringURL("http://example.com:80/about")
	assert.Equal(t, url.String(), "http://example.com/about")
	assert.NoError(t, err)

	// Queries and Fragments currently fail
	// url, err = NormalizeStringURL("http://example.com:80/about?val=true")
	// assert.Equal(t, "http://example.com/about", url.String())
	// assert.NoError(t, err)

	// Not implemented
	url, err = NormalizeStringURL("#")
	assert.Nil(t, url)
	assert.Error(t, err)

	// Relative URL
	url, err = NormalizeStringURL("/about")
	assert.Nil(t, url)
	assert.Error(t, err)
}

func TestNormalizeURL(t *testing.T) {
	surl, _ := url.ParseRequestURI("http://example.com")
	nurl, err := NormalizeURL(surl)
	assert.Equal(t, nurl.String(), "http://example.com")
	assert.NoError(t, err)

	surl, _ = url.ParseRequestURI("http://example.com:80")
	nurl, err = NormalizeURL(surl)
	assert.Equal(t, nurl.String(), "http://example.com")
	assert.NoError(t, err)

	surl, _ = url.ParseRequestURI("http://www.example.com")
	nurl, err = NormalizeURL(surl)
	assert.Equal(t, nurl.String(), "http://www.example.com")
	assert.NoError(t, err)

	// Relative URL
	surl, _ = url.ParseRequestURI("/about")
	nurl, err = NormalizeURL(surl)
	assert.Nil(t, nurl)
	assert.Error(t, err)
}
