package sitemap

import (
	"fmt"

	"github.com/twmb/algoimpl/go/graph"
)

type Sitemapper interface {
	// Add creates a representation of the
	// quad to the Sitemapper
	Add(from string, to string) error

	//SeedURL returns the seed URL of the Sitemap
	SeedURL() string

	//LinksFrom returns the links from a specific node
	LinksFrom(URL string) *[]string
}

type GraphSitemap struct {
	graph    *graph.Graph
	nodemap  map[string]*graph.Node
	hasNodes bool
	root     *graph.Node
}

func NewGraphSitemap() *GraphSitemap {
	nodemap := make(map[string]*graph.Node)
	return &GraphSitemap{
		graph:    graph.New(graph.Directed),
		nodemap:  nodemap,
		hasNodes: false,
	}
}

func (s *GraphSitemap) makeRoot(root *graph.Node) {
	fmt.Printf("Adding ROOT node %s\n", (*root.Value).(string))
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
func (s *GraphSitemap) SeedURL() string {
	return (*s.root.Value).(string)
}

//LinksFrom returns the links from a specific node
// as an unprioritised String slice
func (s *GraphSitemap) LinksFrom(url string) *[]string {
	links := make([]string, 0, 100)
	fmt.Printf("Neighbors of %s:\n", url)
	for _, node := range s.graph.Neighbors(*s.nodemap[url]) {
		val := (*node.Value).(string)
		fmt.Printf("%s\n", val)
		links = append(links, val)
	}
	return &links
}
