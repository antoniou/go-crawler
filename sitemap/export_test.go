// +build integration

package sitemap

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ExportTestSuite struct {
	suite.Suite
}

type MockWriter struct {
	out string
}

func (w *MockWriter) Write(data []byte) (n int, err error) {
	w.out += string(data)
	return len(data), nil
}

func (w *MockWriter) Close() error {
	return fmt.Errorf("Fake error")
}

func (suite *ExportTestSuite) TestExportValidSitemap() {
	seedURL := "http://example.com/"
	s := NewGraphSitemap()
	s.Add(seedURL, seedURL+"contact/")
	s.Add(seedURL, seedURL+"about/")
	s.Add(seedURL, seedURL+"news/")
	s.Add(seedURL+"news/", seedURL)
	s.Add(seedURL+"news/", seedURL+"contact/")
	s.Add(seedURL+"contact/", seedURL+"news/")

	mock := new(MockWriter)
	exp := NewExporter(mock)
	exp.Export(s)

	assert.Equal(suite.T(), strings.TrimSpace(`
	http://example.com/
  http://example.com/contact/
    http://example.com/news/
      http://example.com/
      http://example.com/contact/
  http://example.com/about/
  http://example.com/news/`),
		strings.TrimSpace(mock.out))

}

func (suite *ExportTestSuite) TestFailOnCloseWriter() {
	seedURL := "http://example.com/"
	s := NewGraphSitemap()
	s.Add(seedURL, seedURL+"contact/")

	mock := new(MockWriter)
	exp := NewExporter(mock)
	err := exp.Export(s)

	assert.Error(suite.T(), err)

}

func TestExportTestSuite(t *testing.T) {
	suite.Run(t, new(ExportTestSuite))
}
