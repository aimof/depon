package depon

import (
	"errors"
	"reflect"
	"testing"
)

func sampleAnalyzer() (*analyzer, error) {
	target := "/home/go/src/github.com/aimof/depon/cmd/depon"
	repoRoot := "/home/go/src/github.com/aimof/depon"

	tree, err := newTreeWithRoot(target)
	if err != nil {
		return nil, errors.New("sampleAnalyzer(): cannot initialize sample analyzer")
	}

	return &analyzer{
		tree:     tree,
		target:   target,
		repoRoot: repoRoot,
		srcPath:  "/home/go/src",
	}, nil
}

func TestNewAnalyzer(t *testing.T) {
	ana, err := newAnalyzer(".")
	if err != nil {
		t.Fatalf("%v", err)
	}

	if ana.target != ana.repoRoot {
		t.Error()
	}
}

func TestParseImportedPackages(t *testing.T) {
	s := `github.com/aimof/depon [github.com/aimof/depon/lib]
github.com/aimof/depon/cmd/depon [github.com/aimof/depon fmt github.com/aimof/depon/lib log]
github.com/aimof/depon/lib [encoding/json]`

	want := parseImportedPackages(s)
	got := map[string][]string{
		"github.com/aimof/depon":           []string{"github.com/aimof/depon/lib"},
		"github.com/aimof/depon/cmd/depon": []string{"github.com/aimof/depon", "fmt", "github.com/aimof/depon/lib", "log"},
		"github.com/aimof/depon/lib":       []string{"encoding/json"},
	}

	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:  %v\nwant: %v", got, want)
	}
}

func TestAnalyzeIntoTree(t *testing.T) {
	pkgMap := map[string][]string{
		"github.com/aimof/depon":           []string{"github.com/aimof/depon/lib", "fmt"},
		"github.com/aimof/depon/cmd/depon": []string{"github.com/aimof/depon", "fmt", "github.com/aimof/depon/lib", "log"},
		"github.com/aimof/depon/lib":       []string{"encoding/json", "fmt", "github.com/aimof/depon/lib/foo"},
		"github.com/aimof/depon/lib/foo":   make([]string, 0),
	}

	ana := analyzer{
		target:   "/home/go/src/github.com/aimof/depon/cmd/depon",
		repoRoot: "/home/go/src/github.com/aimof/depon",
		srcPath:  "/home/go/src/",
	}

	var err error
	ana.tree, err = newTreeWithRoot("github.com/aimof/depon/cmd/depon")
	if err != nil {
		t.Error()
	}

	err = ana.analyzeIntoTree("github.com/aimof/depon/cmd/depon", pkgMap)
	if err != nil {
		t.Error()
	}

	wantTree, err := newTreeWithRoot("github.com/aimof/depon/cmd/depon")
	if err != nil {
		t.Error()
	}

	wantTree.root.addChild("github.com/aimof/depon")
	wantTree.root.addChild("fmt")
	wantTree.root.addChild("github.com/aimof/depon/lib")
	wantTree.root.addChild("log")
	child0, ok := wantTree.allNodes["github.com/aimof/depon"]
	if !ok {
		t.Error()
	}
	child0.addChild("github.com/aimof/depon/lib")
	child0.addChild("fmt")
	child1, ok := wantTree.allNodes["github.com/aimof/depon/lib"]
	if !ok {
		t.Error()
	}
	child1.addChild("encoding/json")
	child1.addChild("fmt")
	child1.addChild("github.com/aimof/depon/lib/foo")

	if !reflect.DeepEqual(ana.tree, wantTree) {
		t.Errorf("\ngot:  %v\nwant: %v", ana.tree, wantTree)
	}
}
