package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var nodes = []*MessageNode{
	NewMessageNode("parent", "", "im father", "hello world", nil),
	NewMessageNode("child1", "parent", "im child1", "hello im child1", nil),
	NewMessageNode("child2", "parent", "im child2", "hello im child2", nil),
	NewMessageNode("child3", "parent", "im child3", "hello im child3", nil),

	NewMessageNode(
		"child1child1",
		"child1",
		"im child1 from child1",
		"hello im child1",
		nil,
	),
	NewMessageNode(
		"child1child2",
		"child1",
		"im child2 from child2",
		"hello im child2",
		nil,
	),

	NewMessageNode(
		"child2child",
		"child2",
		"im child1 from child1",
		"odies",
		nil,
	),

	NewMessageNode(
		"child2grandChild",
		"child2child",
		"im grandChild from child2",
		"marvin",
		nil,
	),

	NewMessageNode(
		"child2grandChild2",
		"child2child",
		"im grandChild from child2",
		"coco",
		nil,
	),
	NewMessageNode(
		"child2GreatGrandson",
		"child2grandChild2",
		"im grreate grandson from child2grandChild2",
		"marvin",
		nil,
	),
}

func assertNodes(assert *assert.Assertions, expectedIds []string, children []*MessageNode) {
	assert.Equal(len(expectedIds), len(children))
	for _, child := range children {
		assert.Contains(expectedIds, child.message.id, "Value should be present in the list")
	}
}

func collectChildren(node *MessageNode) []*MessageNode {
	children := make([]*MessageNode, 0)
	node.TransverseInChildren(func(child *MessageNode) {
		children = append(children, child)
	})
	return children
}

func TestNodeInsert(t *testing.T) {
	assert := assert.New(t)
	root := nodes[0]
	others := nodes[1:]

	for _, node := range others {
		root.insert(node)
	}

	type testCase struct {
		desc        string
		node        *MessageNode
		expectedIds []string
	}

	for _, c := range []testCase{
		{
			desc:        "rootChildren",
			node:        root,
			expectedIds: []string{"child1", "child2", "child3"},
		},
		{
			desc:        "Children from child 1",
			node:        others[0],
			expectedIds: []string{"child1child1", "child1child2"},
		},
		{
			desc:        "Children from child 2",
			node:        others[1],
			expectedIds: []string{"child2child"},
		},
		{
			desc:        "grand children from child 2",
			node:        others[5],
			expectedIds: []string{"child2grandChild", "child2grandChild2"},
		},
		{
			desc:        "great children from child 2",
			node:        others[7],
			expectedIds: []string{"child2GreatGrandson"},
		},
	} {
		t.Run(c.desc, func(t *testing.T) {
			children := collectChildren(c.node)
			assertNodes(assert, c.expectedIds, children)
		})
	}

}
