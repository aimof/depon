package depon

import (
	"reflect"
	"testing"
)

func sampleTree() (*Tree, error) {
	pkgMap := map[string][]string{
		"github.com/aimof/depon":           []string{"github.com/aimof/depon/lib", "fmt"},
		"github.com/aimof/depon/cmd/depon": []string{"github.com/aimof/depon", "fmt", "github.com/aimof/depon/lib", "log"},
		"github.com/aimof/depon/lib":       []string{"encoding/json", "fmt", "github.com/aimof/depon/lib/foo"},
		"github.com/aimof/depon/lib/foo":   make([]string, 0),
	}

	ana := analyzer{
		target:   "github.com/aimof/depon/cmd/depon",
		repoRoot: "github.com/aimof/depon",
		srcPath:  "github.com/aimof/",
	}

	var err error
	ana.tree, err = newTreeWithRoot("github.com/aimof/depon/cmd/depon")
	if err != nil {
		return nil, err
	}

	err = ana.analyzeIntoTree("github.com/aimof/depon/cmd/depon", pkgMap)
	if err != nil {
		return nil, err
	}

	return ana.tree, nil
}

func (f Formatter) TestCountAll(t *testing.T) {
	var err error
	f.tree, err = sampleTree()
	if err != nil {
		t.Error()
	}

	want := map[string]CountFormat{
		"github.com/aimof/depon/cmd/depon": {Parents: 0, Children: 2},
		"github.com/aimof/depon":           {Parents: 1, Children: 4},
		"github.com/aimof/depon/lib":       {Parents: 2, Children: 3},
		"github.com/aimof/depon/lib/foo":   {Parents: 1, Children: 0},
		"fmt":                              {Parents: 3, Children: 0},
		"log":                              {Parents: 1, Children: 0},
	}

	got := f.CountAll()
	if err != nil {
		t.Error()
	}

	if !reflect.DeepEqual(got, want) {
		t.Error()
	}
}
