package sitemap

import "github.com/twmb/algoimpl/go/graph"

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

func New() *GraphSitemap {
	nodemap := make(map[string]*graph.Node)
	return &GraphSitemap{
		graph:   graph.New(graph.Directed),
		nodemap: nodemap,
	}
}

func (s *GraphSitemap) addRoot(root string) {
	node := s.addNode(root)
	s.hasNodes = true
	s.root = &node
}

func (s *GraphSitemap) addNode(nodeURL string) graph.Node {
	node := s.graph.MakeNode()
	*node.Value = nodeURL
	s.nodemap[nodeURL] = &node
	return node
}

func (s *GraphSitemap) Add(from string, to string) error {
	if !s.hasNodes {
		s.addRoot(to)
		return nil
	}

	nodeTo := s.addNode(to)

	var nodeFrom *graph.Node
	nodeFrom, ok := s.nodemap[from]
	if !ok {
		nodef := s.addNode(from)
		nodeFrom = &nodef
	}

	s.graph.MakeEdge(*nodeFrom, nodeTo)
	return nil
}

//SeedURL Returns the seed URL (Root) of the Sitemap
func (s *GraphSitemap) SeedURL() string {
	return (*s.root.Value).(string)
}

//LinksFrom returns the links from a specific node
// as an unprioritised String slice
func (s *GraphSitemap) LinksFrom(url string) *[]string {
	links := make([]string, 0, 100)
	for _, node := range s.graph.Neighbors(*s.nodemap[url]) {
		val := (*node.Value).(string)
		links = append(links, val)
	}
	return &links
}
