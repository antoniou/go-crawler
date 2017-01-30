package sitemap

import (
	"fmt"

	"github.com/twmb/algoimpl/go/graph"
)

type Sitemapper interface {
	// Add creates a representation of the
	// quad to the Sitemapper
	Add(quad string) error
}

type GraphSitemap struct {
	graph *graph.Graph
}

func New() *GraphSitemap {
	return &GraphSitemap{
		graph: graph.New(graph.Directed),
	}
}

func (s *GraphSitemap) Add(quad string) error {
	fmt.Printf("Adding node for %s\n", quad)
	node := s.graph.MakeNode()
	*node.Value = quad

	return nil
}
