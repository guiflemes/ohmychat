package telegram

import "fmt"

type messageNode struct {
	firstChild  *messageNode
	nextSibling *messageNode
	id          string
	parent      string
	content     string
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

func (n *messageNode) searchOneLevel(id string) *messageNode {
	if n.id == id {
		return n
	}
	return n.searchChild(id)
}

func (n *messageNode) searchChild(id string) *messageNode {
	if n.firstChild == nil {
		return nil
	}

	child := n.firstChild

	if child.id == id {
		return child
	}

	for child.nextSibling != nil {
		child = child.nextSibling
		if child.id == id {
			return child
		}
	}

	return nil
}

func (n *messageNode) transverseInChildren(fn func(child *messageNode)) {
	if n.firstChild == nil {
		return
	}

	child := n.firstChild
	if child.nextSibling != nil {
		for child != nil {
			fn(child)
			child = child.nextSibling
		}
		return
	}

	fn(child)

}

func (n *messageNode) repChildren() string {
	rep := ""
	count := 1
	n.transverseInChildren(func(child *messageNode) {
		rep += fmt.Sprintf("%d: %s\n", count, child.id)
		count++
	})
	return rep
}

type MessageTree struct {
	root *messageNode
}

func (t *MessageTree) Insert(node *messageNode) *MessageTree {
	if t.root == nil {
		t.root = node
		return t
	}

	t.root.insert(node)
	return t
}

func (t *MessageTree) Search(id string) *messageNode {
	if t.root == nil {
		return nil
	}

	return t.root.searchChild(id)
}

func Fn() *MessageTree {
	tree := &MessageTree{}
	tree.Insert(&messageNode{parent: "", id: "coco"}).
		Insert(&messageNode{parent: "coco", id: "xixi"}).
		Insert(&messageNode{parent: "coco", id: "veve"}).
		Insert(&messageNode{parent: "coco", id: "lolo"}).
		Insert(&messageNode{parent: "xixi", id: "dudu"}).
		Insert(&messageNode{parent: "xixi", id: "caca"}).
		Insert(&messageNode{parent: "lolo", id: "didi"})

	//tree.root.firstChild.printChildren()

	//tree.root.firstChild.nextSibling.nextSibling.printChildren()
	//fmt.Println(tree.root.firstChild.nextSibling.nextSibling.repChildren())
	//fmt.Println(tree.root.firstChild.searchChild("caca"))
	return tree
}

type Engine struct {
	tree *MessageTree
	node *messageNode
}

func NewEngine() *Engine {
	tree := Fn()
	return &Engine{
		tree: tree,
		node: tree.root,
	}
}
func (e *Engine) Reply(message_id string) string {
	node := e.node.searchOneLevel(message_id)
	e.node = node
	return node.content
}
