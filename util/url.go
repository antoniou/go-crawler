package util

import (
	"net/url"

	"github.com/goware/urlx"
)

// NormalizeURL normalizes a url.URL to its canonical form
// and returns a url.URL
func NormalizeURL(orig *url.URL) (*url.URL, error) {
	normURL, err := urlx.Normalize(orig)
	if err != nil {
		return nil, err
	}
	u, err := url.ParseRequestURI(normURL)
	return u, err
}

// NormalizeStringURL normalizes a String URL to its canonical form
// and returns a url.URL
func NormalizeStringURL(orig string) (*url.URL, error) {
	normURL, err := urlx.NormalizeString(orig)
	if err != nil {
		return nil, err
	}
	return url.ParseRequestURI(normURL)
}
