package depon

import (
	"os"
	"os/exec"
	"strings"
)

// Formatter outputs formatted text
type Formatter struct {
	tree *Tree
}

// NewFormatter returns Formatter which contains analyzed Tree.
func NewFormatter(s string) (*Formatter, error) {
	ana, err := newAnalyzer(s)
	if err != nil {
		return nil, err
	}

	err = os.Chdir(ana.repoRoot)
	if err != nil {
		return nil, err
	}

	b, err := exec.Command("/bin/sh", "-c", `go list -f "{{.ImportPath}} {{.Imports}}" ./...`).Output()
	if err != nil {
		return nil, err
	}

	pkgMap := parseImportedPackages(strings.TrimSuffix(string(b), "\n"))

	err = ana.analyzeIntoTree(strings.TrimPrefix(ana.target, ana.srcPath), pkgMap)
	if err != nil {
		return nil, err
	}

	return ana.ToFormatter(), nil
}

// CountAll counts all node's parentsa and children number.
func (f Formatter) CountAll() map[string]CountFormat {
	countMap := make(map[string]CountFormat, len(f.tree.allNodes))

	for key, value := range f.tree.allNodes {
		countMap[key] = CountFormat{Parents: len(value.parents), Children: len(value.children)}
	}
	return countMap
}
