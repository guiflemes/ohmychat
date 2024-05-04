package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

var stubNodes = func() []*MessageNode {
	return []*MessageNode{
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
			"im great grandson from child2grandChild2",
			"marvin",
			nil,
		),
	}
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
	nodes := stubNodes()
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

func TestNodeSearchChild(t *testing.T) {
	assert := assert.New(t)
	nodes := stubNodes()
	root := nodes[0]
	others := nodes[1:]

	for _, node := range others {
		root.insert(node)
	}

	type testCase struct {
		desc          string
		node          *MessageNode
		expectedChild *MessageNode
		searchID      string
	}

	for _, c := range []testCase{
		{
			desc:          "search for root's child",
			node:          root,
			expectedChild: others[0],
			searchID:      others[0].message.id,
		},
		{
			desc:          "search for child1's child",
			node:          others[0],
			expectedChild: others[3],
			searchID:      others[3].message.id,
		},
		{
			desc:          "search for child2grandChild2's child",
			node:          others[7],
			expectedChild: others[8],
			searchID:      others[8].message.id,
		},
		{
			desc:          "search for child3's child",
			node:          others[1],
			expectedChild: nil,
			searchID:      "someID",
		},
	} {
		t.Run(c.desc, func(t *testing.T) {
			n := c.node.searchChild(c.searchID)
			assert.Equal(n, c.expectedChild)
		})

	}

}

func TestNodeTransverseInChildren(t *testing.T) {
	assert := assert.New(t)
	nodes := stubNodes()
	root := nodes[0]
	others := nodes[1:]

	for _, node := range others {
		root.insert(node)
	}

	type testCase struct {
		desc           string
		node           *MessageNode
		expectedResult int
	}

	for _, c := range []testCase{
		{
			desc:           "root node contains 3 children",
			node:           root,
			expectedResult: 3,
		},
		{
			desc:           "child1 node contains 2 children",
			node:           others[0],
			expectedResult: 2,
		},
		{
			desc:           "child2grandChild2 node contains 1 child",
			node:           others[7],
			expectedResult: 1,
		},
		{
			desc:           "child2grandChild2 node doest have any node",
			node:           others[8],
			expectedResult: 0,
		},
	} {
		t.Run(c.desc, func(t *testing.T) {
			var countChildren int

			c.node.TransverseInChildren(func(child *MessageNode) {
				countChildren++

			})

			assert.Equal(countChildren, c.expectedResult)

		})
	}

}

type MessageTreeSuite struct {
	suite.Suite
	n []*MessageNode
}

func (m *MessageTreeSuite) BeforeTest(suiteName, testName string) {
	m.n = stubNodes()
}

func (m *MessageTreeSuite) TestSetRoot() {
	tree := &MessageTree{}

	for _, node := range m.n {
		m.Run(fmt.Sprintf("set root to node %s", node.message.id), func() {
			tree.SetRoot(node)
			m.Equal(tree.Root(), node)
		})
	}
}

func (m *MessageTreeSuite) TestInsert() {
	tree := &MessageTree{}

	for _, node := range m.n {
		m.Run(fmt.Sprintf("insert node %s", node.message.id), func() {
			root := tree.Insert(node)
			m.Equal(tree.Root().Message().ID(), root.Root().Message().ID())
		})
	}
}

func (m *MessageTreeSuite) TestSearch() {
	tree := &MessageTree{}

	for _, node := range m.n {
		tree.Insert(node)
	}

	n := tree.Search("child1")
	m.Equal(n.Message().ID(), "child1")

	n = tree.Search("child2")
	m.Equal(n.Message().ID(), "child2")

}

func TestMessageTreeSuite(t *testing.T) {
	suite.Run(t, new(MessageTreeSuite))
}
