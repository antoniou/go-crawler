package sitemap

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ExportTestSuite struct {
	suite.Suite
}

func (suite *ExportTestSuite) TestStopFetcher() {

}

func (suite *ExportTestSuite) TestNewAsyncHTTPFetcherConstructor() {
	f := NewAsyncHTTPFetcher()
	assert.Implements(suite.T(), (*Fetcher)(nil), f)
	assert.NotNil(suite.T(), f)
}

func TestExportTestSuite(t *testing.T) {
	suite.Run(t, new(ExportTestSuite))
}
