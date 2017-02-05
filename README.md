# go-crawler [![Build Status](https://travis-ci.org/antoniou/go-crawler.svg?branch=master)](https://travis-ci.org/antoniou/go-crawler)

A Web Crawler written in Go

The crawler is limited to one domain - Given domain http://www.example.com it will crawl all pages within the domain, but will not follow external links.

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
$ go-crawler -v -o tom_sitemap.out http://tomblomfield.com
```

## Assumptions
Certain assumptions about the requirements should be made:

1. The crawler does crawl subdomains of the specific domain provided
1. If a specific path within a domain is given


Concerns/Improvements:
1. Parallelize implementation even further
1. Bloom Filter False-positive Rate:
1. URLs that time-out
1. Use robots.txt and existing sitemap
1. At the moment, the crawler limits itself within a single scheme (e.g, http). For example, if  http://www.example.com is given as input, the crawler will not follow links to https://www.example.com/about. This is an improvement that the crawler needs


### Asymptotic Complexity:
Space: O(n)
Time: O(n)
