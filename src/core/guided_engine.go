package core

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"oh-my-chat/src/logger"
	"oh-my-chat/src/models"
)

type Action interface {
	Handle(ctx context.Context, message *models.Message) error
}

type Message struct {
	id      string
	parent  string
	Content string
	Action  Action
}

func (m Message) ID() string {
	return m.id
}

func (m Message) HasAction() bool {
	return m.Action != nil
}

type MessageNode struct {
	firstChild  *MessageNode
	nextSibling *MessageNode
	message     Message
}

func (n *MessageNode) Message() Message {
	return n.message
}

func (n *MessageNode) insert(node *MessageNode) {
	if n.message.id == node.message.parent {
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

	if n.firstChild.message.id == node.message.parent {
		n.firstChild.insert(node)
		return
	}

	found := false
	sibling := n.firstChild
	for sibling.nextSibling != nil {
		sibling = sibling.nextSibling
		if sibling.message.id == node.message.parent {
			found = !found
			break
		}
	}

	if !found {
		fmt.Printf(
			"node %s without parent, given parent %s not found/n",
			node.message.id,
			node.message.parent,
		)
		return
	}

	sibling.insert(node)

}

func (n *MessageNode) SearchOneLevel(id string) *MessageNode {
	if n.message.id == id {
		return n
	}
	return n.searchChild(id)
}

func (n *MessageNode) searchChild(id string) *MessageNode {
	if n.firstChild == nil {
		return nil
	}

	child := n.firstChild

	if child.message.id == id {
		return child
	}

	for child.nextSibling != nil {
		child = child.nextSibling
		if child.message.id == id {
			return child
		}
	}

	return nil
}

func (n *MessageNode) TransverseInChildren(fn func(child *MessageNode)) {
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

func (n *MessageNode) RepChildren() string {
	rep := ""
	count := 1
	n.TransverseInChildren(func(child *MessageNode) {
		rep += fmt.Sprintf("%d: %s\n", count, child.message.id)
		count++
	})
	return rep
}

type MessageTree struct {
	root *MessageNode
}

func (t *MessageTree) Root() *MessageNode {
	return t.root
}

func (t *MessageTree) SetRoot(root *MessageNode) {
	t.root = root
}

func (t *MessageTree) Insert(node *MessageNode) *MessageTree {
	if t.root == nil {
		t.root = node
		return t
	}

	t.root.insert(node)
	return t
}

func (t *MessageTree) Search(id string) *MessageNode {
	if t.root == nil {
		return nil
	}

	return t.root.searchChild(id)
}

type GuidedResponseEngine struct {
	tree         *MessageTree
	node         *MessageNode
	dialogLaunch bool
	actionQueue  ActionQueue
	setup        bool
}

func (e *GuidedResponseEngine) Config(workflow Workflow) {}

func (e *GuidedResponseEngine) resolveMessageNode(messageID string) {

	if e.dialogLaunch {
		node := e.node.SearchOneLevel(messageID)

		if node == nil {
			e.node = e.tree.Root()
			e.dialogLaunch = false
			return
		}
		e.node = node
		return
	}

	e.dialogLaunch = true

}

func (e *GuidedResponseEngine) GetActionQueue() ActionQueue {
	return e.actionQueue
}

func (e *GuidedResponseEngine) Name() string {
	return "guided"
}

func (e *GuidedResponseEngine) HandleMessage(input models.Message, output chan<- models.Message) {
	// TODO fix context
	ctx := context.Background()

	if !e.setup {
		logger.Logger.Error("engine is not ready", zap.String("context", "guided_engine"))
		response := &input
		response.Output = "some error ocurred, please contant admin"
		output <- *response
		return
	}

	e.resolveMessageNode(input.ID)

	if e.node.Message().HasAction() {
		actionPair := ActionReplyPair{replyTo: output, action: e.node.message.Action, input: input}
		queue := e.GetActionQueue()
		queue.Put(ctx, actionPair)
		return
	}

	response := &input
	response.Output = "SomeMessage"
	output <- *response
}
