package depon

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewTree(t *testing.T) {
	tree := newTree()
	if tree.root != nil {
		t.Error()
	}
	if tree.allNodes == nil {
		t.Error()
	}
	if len(tree.allNodes) != 0 {
		t.Error()
	}
}

func TestMNAdd(t *testing.T) {
	tree := newTree()

	mn := newMapNodes()
	name := "name"
	nodeptr := tree.newNode("name")
	err := mn.add(nodeptr)
	if err != nil {
		t.Error()
	}
	if v, ok := mn[nodeptr.name]; !ok {
		t.Error()
	} else if v.name != name {
		t.Error()
	}
	err = mn.add(nodeptr)
	if err != nil {
		t.Error()
	}

	var nilptr *node
	err = mn.add(nilptr)
	if err == nil {
		t.Error()
	}

	dummyptr := tree.newNode("name")
	err = mn.add(dummyptr)
	if err == nil {
		t.Error()
	}
}

func TestSetRoot(t *testing.T) {
	tree := newTree()
	name := "root"
	err := tree.setRoot(name)
	if err != nil {
		t.Error()
	}
	if v, ok := tree.allNodes[name]; !ok {
		t.Error()
	} else if v.name != name {
		t.Error()
	}
}

func TestNewNode(t *testing.T) {
	tree := newTree()

	name := "name"
	node := tree.newNode(name)
	if node.name != name {
		t.Error()
	}
	if node.tree != tree {
		t.Error()
	}
	if node.parents == nil || len(node.parents) != 0 {
		t.Error()
	}
	if node.children == nil || len(node.children) != 0 {
		t.Error()
	}
}

func TestNewTreeWithRoot(t *testing.T) {
	name := "root"
	tree, err := newTreeWithRoot(name)
	if err != nil {
		t.Error()
	}
	if tree.root.name != name {
		t.Error()
	}
	if len(tree.allNodes) != 1 {
		t.Error()
	}
}

func makeSampleTree() (*Tree, error) {
	tree, err := newTreeWithRoot("root")
	if err != nil {
		return nil, errors.New("")
	}

	tree.root.addChild("child0")
	tree.root.addChild("child1")
	for _, v := range tree.root.children {
		v.addChild("grandson")
	}
	return tree, nil
}

func TestCount(t *testing.T) {
	tree, err := makeSampleTree()
	if err != nil {
		t.Error()
	}

	got, err := tree.count("child1")
	if err != nil {
		t.Error()
	}
	want := CountFormat{
		Parents:  1,
		Children: 1,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:  %v\nwant: %v", got, want)
	}

	got, err = tree.count("grandson")
	if err != nil {
		t.Error()
	}
	want = CountFormat{
		Parents:  2,
		Children: 0,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot:  %v\nwant: %v", got, want)
	}
}

func TestAddChild(t *testing.T) {
	rootName := "root"
	tree, err := newTreeWithRoot(rootName)
	if err != nil {
		t.Error()
	}
	childName := "child"
	err = tree.root.addChild(childName)
	if err != nil {
		t.Error()
	}
	if child, ok := tree.root.children[childName]; !ok {
		t.Error()
	} else if child == nil {
		t.Error()
	}

	if child, ok := tree.allNodes[childName]; !ok {
		t.Error()
	} else if _, ok := child.parents[tree.root.name]; !ok {
		t.Error()
	}

	err = tree.root.addChild(childName)
	if err != nil {
		t.Error()
	}
}
