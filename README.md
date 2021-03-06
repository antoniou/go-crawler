# go-crawler [![Build Status](https://travis-ci.org/antoniou/go-crawler.svg?branch=master)](https://travis-ci.org/antoniou/go-crawler)

A Web Crawler written in Go. This was developed with the purpose of me trying to learn to write idiomatic Go

The crawler is limited to one domain - Given domain http://www.example.com it will crawl all pages within the domain, but will not follow external links.

The crawler will also provide link information between the pages of the crawled domain. For example, given a link from http://example.com/ to http://example.com/about/ and another link from http://example.com/about/ to http://example.com/ the output will be:

```text
http://example.com/
  http://example.com/about/
    http://example.com/
```

## Installation
To install go-crawler, you'll need to have Golang installed and environment variable [$GOPATH appropriately set](https://golang.org/doc/install).
```bash
$ go get github.com/antoniou/go-crawler
```

## Usage
```bash
$ go-crawler -o tom_sitemap.out http://tomblomfield.com
```

To see what happens during crawling, enable verbose mode:
```bash
$ go-crawler --verbose -o tom_sitemap.out http://tomblomfield.com
```

## Assumptions
Certain assumptions about the requirements should be made:

1. The crawler does not crawl subdomains of the specific domain provided.
1. If a specific path within a domain is given as a seed URL ( e.g http://tomblomfield.com/about), then crawling will only happen within that path.
1. Crawling happens only for a specific scheme (e.g, http or https, not both)


## Architecture
The solution goes through two major phases:

1. Crawling the site and creating an in-memory representation of a sitemap.
1. Exporting the sitemap to a file representation.

### Crawling the site
Crawling happens with asynchronous channel-based communication between the following components:

1. Fetcher: Awaits for requests to fetch pages, and hands over responses to the requests to the Parser
2. Parser: Awaits for http responses (from Fetcher), parses the responses and hands over URLs that are found to the Tracker
3. Tracker: Awaits for URLs that have been found (from Parser) and checks whether the URLs have been already crawled. If not, the Tracker hands over new requests to the fetcher.
4. Sitemapper: Holds the sitemap representation and awaits to receive new nodes and edges to add to the sitemap (from the Tracker)

The Crawler is the orchestrating component that starts all the workers and awaits for all the workers to be in "WAITING" state, to detect that the work has finished.

![crawl](https://raw.githubusercontent.com/antoniou/go-crawler/master/dotgraph/crawlGraph.png "Crawling stage architecture")

### Exporting the sitemap
After crawling the site, the sitemap is exported using an Exported. The default Exporter outputs the sitemap to a file using an indented structure to represent the links between pages.
The Crawler initiates the export by creating a new Exporter. The exporter then queries the Sitemapper to walk through the Graph

![export](https://github.com/antoniou/go-crawler/raw/master/dotgraph/exportgraph.png "Exporting sitemap stage architecture")

#### Average Space Complexity :
The solution makes use of bloom filters, graphs and hashmaps:
Given N crawled pages and M links between the pages, each of average size L, their space complexity is:

1. Bloom Filters used for pages: O(1) - fixed Space
2. HashMap used for pages: O(N)
3. Graph Nodes used for pages: O(N)
4. Graph Edges used for links: O(M)
5. Storing the pages to be parsed: O(L), since the current Parser implementation parses one page at a time.

Therefore, the average space complexity is linear to the maximum of pages and links between them:
```
O(N + M + L)
```

## Performance
At the moment there are 4 asynchronous workers that handle the crawling. Components are working in parallel and are communicating asynchronously, however each component is currently single threaded.

To improve performance, we can have multiple workers of each type waiting to receive and process work. To achieve this, we'd need to extract the channels from each individual component, so that each channel can have multiple receivers.

When crawling large documents, performance can be improved by chunking the document between several Parsers.

We can also take advantage of files like robots.txt and the provided website sitemap as hints to the website structure.

## Future Work/Improvements:
1. Parallelize implementation even further as described in [Performance](#Performance)
1. Bloom Filter False-positive Rate: We make use of [Bloom Filters](https://en.wikipedia.org/wiki/Bloom_filter?oldformat=true) to detect whether a page has been crawled or not. If the size of the bloom filters is not adequately large, it can have a false positive response. To deal with this, we'd have to create an adequately large filter.  
1. Improve coverage of unit tests.
1. At the moment, a site is completely crawled into a Sitemap in memory and then exported/shown. This will not work for large websites whose sitemap might not even fit in memory.
1. At the moment, the crawler limits itself within a single scheme (e.g, http). For example, if  http://www.example.com is given as input, the crawler will not follow links to https://www.example.com/about. This is an improvement that the crawler needs.
1. Queries (e.g, ?account=21312) and Fragments (#About) are not removed, hence a page can be crawled multiple times because of them.
