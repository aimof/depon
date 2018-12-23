package depon

import (
	"errors"
	"fmt"
)

// Tree is a Package tree.
type Tree struct {
	root     *node
	allNodes mapNodes
}

// mapNodes is a map of Nodes. keys are Node name.
type mapNodes map[string]*node

// Node is a node info.
type node struct {
	name     string
	tree     *Tree
	parents  mapNodes
	children mapNodes
}

func newTreeWithRoot(name string) (*Tree, error) {
	tree := newTree()
	err := tree.setRoot(name)
	return tree, err
}

// NewTree returns initialized tree with a simgle nil node.
func newTree() *Tree {
	return &Tree{
		root:     nil,
		allNodes: newMapNodes(),
	}
}

func (tree *Tree) setRoot(packageName string) error {
	root := tree.newNode(packageName)
	tree.root = root
	return tree.allNodes.add(tree.root)
}

// CountFormat is count of parents and children of a node.
type CountFormat struct {
	Parents  int
	Children int
}

func (tree *Tree) count(name string) (CountFormat, error) {
	targetNode, ok := tree.allNodes[name]
	if !ok || targetNode == nil {
		return CountFormat{}, fmt.Errorf("count: %s doesn't exist", name)
	}

	return CountFormat{
		Parents:  len(targetNode.parents),
		Children: len(targetNode.children),
	}, nil
}

func newMapNodes() mapNodes {
	return make(mapNodes, 64)
}

func (mn mapNodes) add(nodeptr *node) error {
	if nodeptr == nil {
		return errors.New("mapNodes.add(): nodeptr is nil")
	}
	if v, ok := mn[nodeptr.name]; ok {
		if v != nodeptr {
			return errors.New("mapNodes.add(): nodeptr is not same as ptr already resisterd")
		}
		return nil
	}
	mn[nodeptr.name] = nodeptr
	return nil
}

// NewNode returns a node without children.
func (tree *Tree) newNode(name string) *node {
	return &node{
		name:     name,
		tree:     tree,
		parents:  newMapNodes(),
		children: newMapNodes(),
	}
}

func (pn *node) addChild(name string) error {
	if existedChildNode, ok := pn.tree.allNodes[name]; ok {
		pn.children.add(existedChildNode)
		existedChildNode.parents.add(pn)
		return nil
	}

	childNode := pn.tree.newNode(name)
	pn.tree.allNodes.add(childNode)

	err := pn.children.add(childNode)
	if err != nil {
		return err
	}
	err = childNode.parents.add(pn)
	if err != nil {
		return err
	}
	return nil
}
