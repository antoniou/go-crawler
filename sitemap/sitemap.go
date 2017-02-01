package sitemap

import (
	"fmt"

	"github.com/twmb/algoimpl/go/graph"
)

type Sitemapper interface {
	// Add creates a representation of the
	// quad to the Sitemapper
	Add(from string, to string) error
}

type GraphSitemap struct {
	graph    *graph.Graph
	nodemap  map[string]*graph.Node
	hasNodes bool
	root     *graph.Node
}

func New() *GraphSitemap {
	nodemap := make(map[string]*graph.Node)
	return &GraphSitemap{
		graph:   graph.New(graph.Directed),
		nodemap: nodemap,
	}
}

func (s *GraphSitemap) Add(from string, to string) error {
	fmt.Printf("1.Adding node for %s\n", to)
	node := s.graph.MakeNode()
	*node.Value = to
	s.nodemap[to] = &node

	_, ok := s.nodemap[from]
	if !ok {
		fmt.Printf("1.Adding node for %s\n", from)
		nodeFrom := s.graph.MakeNode()
		*nodeFrom.Value = from
		s.nodemap[from] = &nodeFrom
		if !s.hasNodes {
			fmt.Printf("Root node is %s\n", from)
			s.hasNodes = true
			s.root = &nodeFrom
		}
	}
	s.graph.MakeEdge(*s.nodemap[from], node)

	return nil
}
