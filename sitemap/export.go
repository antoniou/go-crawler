package sitemap

import (
	"io"

	"github.com/willf/bloom"
)

// Exporter takes a Sitemapper (Sitemap represenation)
// and creates a representation of it
type Exporter interface {
	Export(Sitemapper) error
}

// FileExporter exports a Sitemapper to a File
type FileExporter struct {
	writer io.WriteCloser
	filter *bloom.BloomFilter
}

// Export exports Sitemapper s to FileExporter.writer
// Returns nil or error on failure
func (f *FileExporter) Export(s Sitemapper) error {
	err := f.exportRecursive(s, s.SeedURL(), "")
	if err != nil {
		return err
	}

	f.writer.Close()
	return nil
}

func (f *FileExporter) exportRecursive(s Sitemapper, node string, indentation string) error {
	_, err := f.writer.Write([]byte(indentation + node + "\n"))
	if err != nil {
		return err
	}

	if !f.filter.TestAndAddString(node) {
		ind := indentation + "  "
		links := *s.LinksFrom(s.SeedURL())
		for _, link := range links {
			f.exportRecursive(s, link, ind)
		}
	}

	return nil
}

// NewExporter is an Exporter constructor
func NewExporter(w io.WriteCloser) *FileExporter {
	filter := bloom.New(20000, 5)
	return &FileExporter{
		writer: w,
		filter: filter,
	}
}
