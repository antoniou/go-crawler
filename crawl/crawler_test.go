// +build integration

package crawl

import (
	"testing"

	"github.com/antoniou/go-crawler/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type CrawlTestSuite struct {
	suite.Suite
}

func (suite *CrawlTestSuite) TestInvalidInputCrawler() {
	seedURL, _ := util.NormalizeStringURL("http://notExistingUrl404.com")
	crawler := NewAsyncHTTPCrawler(seedURL)
	sitemap, err := crawler.Crawl()

	assert.NoError(suite.T(), err)
	_, err = sitemap.SeedURL()
	assert.Error(suite.T(), err)

	seedURL, _ = util.NormalizeStringURL("ftp://invalidscheme.com")
	crawler = NewAsyncHTTPCrawler(seedURL)
	sitemap, err = crawler.Crawl()
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), sitemap)
}

func (suite *CrawlTestSuite) TestValidInputCrawler() {
	seedURL, _ := util.NormalizeStringURL("http://tomblomfield.com/about")
	crawler := NewAsyncHTTPCrawler(seedURL)
	sitemap, err := crawler.Crawl()
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), sitemap)
	seed, _ := sitemap.SeedURL()
	assert.Equal(suite.T(), "http://tomblomfield.com/about", seed)

}

func TestCrawlTestSuite(t *testing.T) {
	suite.Run(t, new(CrawlTestSuite))
}
