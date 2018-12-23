package depon

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type analyzer struct {
	tree     *Tree
	target   string
	repoRoot string
	srcPath  string
}

func newAnalyzer(targetPath string) (*analyzer, error) {
	gopath := os.Getenv("GOPATH")

	targetAbs, err := filepath.Abs(targetPath)
	if err != nil {
		return nil, err
	}

	err = os.Chdir(targetAbs)
	if err != nil {
		return nil, err
	}
	b, err := exec.Command("/bin/sh", "-c", "git rev-parse --show-toplevel").Output()
	if err != nil {
		return nil, err
	}
	repoRoot, err := filepath.Abs(strings.TrimSuffix(string(b), "\n"))

	if err != nil {
		return nil, err
	}

	if !strings.HasPrefix(targetAbs, repoRoot) {
		return nil, fmt.Errorf("newAnalyzer\ntarget: %s\nrepo:   %s", targetAbs, repoRoot)
	}

	tree, err := newTreeWithRoot(strings.TrimPrefix(targetAbs, gopath+"/src/"))
	if err != nil {
		return nil, err
	}

	return &analyzer{
		tree:     tree,
		target:   targetAbs,
		repoRoot: repoRoot,
		srcPath:  gopath + "/src/",
	}, nil
}

func (ana *analyzer) ToFormatter() *Formatter {
	return &Formatter{
		tree: ana.tree,
	}
}

func parseImportedPackages(s string) map[string][]string {
	rows := strings.Split(s, "\n")

	importedPkgs := make(map[string][]string, 0)
	for _, row := range rows {
		if len(row) < 3 {
			continue
		}
		words := strings.Split(row, " ")
		packages := make([]string, 0, len(words))
		for _, word := range words {
			if word != "" {
				packages = append(packages, word)
			}
		}

		switch len(words) {
		case 0:
			continue
		case 1:
			importedPkgs[words[0]] = make([]string, 0)
		default:
			words[1] = strings.TrimPrefix(words[1], "[")
			words[len(words)-1] = strings.TrimSuffix(words[len(words)-1], "]")
			importedPkgs[words[0]] = words[1:]
		}
	}
	return importedPkgs
}

func (ana *analyzer) analyzeIntoTree(name string, pkgMap map[string][]string) error {
	prefix := strings.TrimPrefix(ana.repoRoot, ana.srcPath)
	pkgs, ok := pkgMap[name]
	if !ok {
		return fmt.Errorf("analyzeTree: %s does not exist in pkgMap", name)
	}

	n, ok := ana.tree.allNodes[name]
	if !ok {
		return fmt.Errorf("analyzeTree: %s does not exist in tree", name)
	}

	for _, importedPkg := range pkgs {
		n.addChild(importedPkg)
		if strings.HasPrefix(importedPkg, prefix) {
			err := ana.analyzeIntoTree(importedPkg, pkgMap)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
