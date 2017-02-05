# go-crawler [![Build Status](https://travis-ci.org/antoniou/go-crawler.svg?branch=master)](https://travis-ci.org/antoniou/go-crawler)

A Web Crawler written in Go

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

1. The crawler does crawl subdomains of the specific domain provided
1. If a specific path within a domain is given


## Architecture
The solution goes through two major phases:

1. Crawling the site and creating an in-memory representation of a sitemap
1. Exporting the sitemap to a file representation

### Crawling the site
Crawling happens with asynchronous channel-based communication between the following components:
1. Fetcher: Awaits for requests to fetch pages, and hands over responses to the requests to the Parser
2. Parser: Awaits for http responses (from Fetcher), parses the responses and hands over URLs that are found to the Tracker
3. Tracker: Awaits for URLs that have been found (from Parser) and checks whether the URLs have been already crawled. If not, the Tracker hands over new requests to the fetcher

![crawl](https://github.com/antoniou/go-crawler/dotgraph/crawlGraph.png "Crawling stage architecture")

## Asymptotic Complexity:
#### Space Complexity :
The solution makes use of bloom filters, graphs and hashmaps:
Given N crawled pages and M links between the pages, their space complexity is:
1. Bloom Filters used for pages: O(1) - fixed Space
2. HashMap used for pages: O(N)
3. Graph Nodes used for pages: O(N)
3. Graph Edges used for links: O(M)

Therefore, the asymptotic space complexity is linear to the maximum of pages and links between them:
```
O(max(N, M))
```
#### Time Complexity:
The solution goes through the following phases:
1. Crawling the site, and making sure that the same page is not crawled twice
2. Creating an in-memory representation of the site-map
3. Exporting the sitemap to a file representation

Given N crawled pages and M links between the pages, the time complexity of each step should be:
1. Crawling: O(N)
1. In-memory representation: O(max(N, M))
1. Exporting sitemap: O(max(N, M))

The overall time complexity should then be:
```
O(max(N, M))
```



## Future Work/Improvements:
1. Parallelize implementation even further
1. Bloom Filter False-positive Rate:
1. URLs that time-out
1. Use robots.txt and existing sitemap
1. At the moment, the crawler limits itself within a single scheme (e.g, http). For example, if  http://www.example.com is given as input, the crawler will not follow links to https://www.example.com/about. This is an improvement that the crawler needs
