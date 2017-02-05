package sitemap

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type SitemapTestSuite struct {
	suite.Suite
}

func (suite *SitemapTestSuite) TestStopFetcher() {

}

func (suite *SitemapTestSuite) TestNewAsyncHTTPFetcherConstructor() {
}

func TestSitemapTestSuite(t *testing.T) {
	suite.Run(t, new(SitemapTestSuite))
}
