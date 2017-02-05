package sitemap

import (
	"fmt"

	"github.com/antoniou/go-crawler/util"
	"github.com/twmb/algoimpl/go/graph"
)

// A Sitemapper holds the represenation of
// a sitemap. Links between URLs are created
// with Add
type Sitemapper interface {
	// Add creates a representation of a link
	Add(from string, to string) error

	//SeedURL returns the seed URL of the Sitemap
	// or error in case there is none
	SeedURL() (string, error)

	//LinksFrom returns the links from a specific node
	LinksFrom(URL string) *[]string
}

// GraphSitemap is a Directed Graph-based
// implementation of Sitemapper
type GraphSitemap struct {
	graph    *graph.Graph
	nodemap  map[string]*graph.Node
	hasNodes bool
	root     *graph.Node
}

// NewGraphSitemap constructs a GraphSitemap
// It needs to maintain a nodemap:
// url -> graph.Node(url), e.g,
// "http://example.com" -> graph.Node("http://example.com" )
func NewGraphSitemap() *GraphSitemap {
	nodemap := make(map[string]*graph.Node)
	return &GraphSitemap{
		graph:    graph.New(graph.Directed),
		nodemap:  nodemap,
		hasNodes: false,
	}
}

// Add creates a representation of a link
// GraphSitemap creates a Graph Node for the from and to URLs
// It also creates an edge for the link between them
func (s *GraphSitemap) Add(from string, to string) error {
	nodeFrom, _ := s.addNode(from)
	if !s.hasNodes {
		s.makeRoot(nodeFrom)
	}
	nodeTo, _ := s.addNode(to)

	// Add edge between from and to nodes
	return s.graph.MakeEdge(*nodeFrom, *nodeTo)
}

//SeedURL Returns the seed URL (Root) of the Sitemap
func (s *GraphSitemap) SeedURL() (string, error) {
	if s.root == nil {
		return "", fmt.Errorf("Sitemap could not be exported")
	}
	return (*s.root.Value).(string), nil
}

//LinksFrom returns the links from a specific node
// as an unprioritised String slice
func (s *GraphSitemap) LinksFrom(url string) *[]string {
	links := make([]string, 0, 100)
	util.Printf("Neighbors of %s:\n", url)
	for _, node := range s.graph.Neighbors(*s.nodemap[url]) {
		val := (*node.Value).(string)
		util.Printf("%s\n", val)
		links = append(links, val)
	}
	return &links
}

func (s *GraphSitemap) makeRoot(root *graph.Node) {
	util.Printf("Adding ROOT node %s\n", (*root.Value).(string))
	s.hasNodes = true
	s.root = root
}

func (s *GraphSitemap) addNode(nodeURL string) (*graph.Node, error) {
	n, ok := s.nodemap[nodeURL]
	if ok {
		return n, fmt.Errorf("URL Node %s already exists", nodeURL)
	}

	node := s.graph.MakeNode()
	*node.Value = nodeURL
	s.nodemap[nodeURL] = &node
	return &node, nil
}
