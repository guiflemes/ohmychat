package telegram

import "fmt"

type content interface {
	hash() int
}

type messageNode struct {
	firstChild  *messageNode
	nextSibling *messageNode
	id          string
	parent      string
}

//                           coco
//     xixi (coco)                     veve(coco)     lolo(coco)
// dudu(xixi)   caca(xixi)                               didi(lolo)

func (n *messageNode) insert(node *messageNode) {
	if n.id == node.parent {
		if n.firstChild == nil {
			n.firstChild = node
			return
		}
		sibling := n.firstChild
		for sibling.nextSibling != nil {
			sibling = sibling.nextSibling
		}
		sibling.nextSibling = node
		return
	}

	if n.firstChild.id == node.parent {
		n.firstChild.insert(node)
		return
	}

	found := false
	sibling := n.firstChild
	for sibling.nextSibling != nil {
		sibling = sibling.nextSibling
		if sibling.id == node.parent {
			found = !found
			break
		}
	}

	if !found {
		fmt.Printf("node %s without parent, given parent %s not found/n", node.id, node.parent)
		return
	}

	sibling.insert(node)

}

func (n *messageNode) printChildren() {
	str := ""
	count := 1
	if n.firstChild != nil {
		str += fmt.Sprintf("%d: %s\n", count, n.firstChild.id)

		nextSibling := n.firstChild.nextSibling
		for nextSibling != nil {
			count++
			str += fmt.Sprintf("%d: %s\n", count, nextSibling.id)
			nextSibling = nextSibling.nextSibling
		}
	}
	fmt.Println(str)
}

func (n *messageNode) search(value int) *messageNode {
	return n
}

type MessageTree struct {
	root *messageNode
}

func (t *MessageTree) insert(node *messageNode) *MessageTree {
	if t.root == nil {
		t.root = node
		return t
	}

	t.root.insert(node)
	return t
}

func (t *MessageTree) Search(value int) *messageNode {
	if t.root == nil {
		return nil
	}
	return t.root.search(value)
}

func Fn() {
	tree := &MessageTree{}
	tree.insert(&messageNode{parent: "", id: "coco"})
	tree.insert(&messageNode{parent: "coco", id: "xixi"})
	tree.insert(&messageNode{parent: "coco", id: "veve"})
	tree.insert(&messageNode{parent: "coco", id: "lolo"})

	tree.insert(&messageNode{parent: "xixi", id: "dudu"})
	tree.insert(&messageNode{parent: "xixi", id: "caca"})

	tree.insert(&messageNode{parent: "lolo", id: "didi"})
	//tree.root.firstChild.printChildren()

	tree.root.firstChild.nextSibling.nextSibling.printChildren()
}
