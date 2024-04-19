package core

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"oh-my-chat/src/actions"
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

type guidedResponseEngine struct {
	tree         *MessageTree
	node         *MessageNode
	dialogLaunch bool
	actionQueue  ActionQueue
	setup        bool
}

func NewGuidedResponseEngine(actionQueue ActionQueue) *guidedResponseEngine {
	return &guidedResponseEngine{actionQueue: actionQueue}
}

func (e *guidedResponseEngine) IsReady() bool {
	return e.setup
}

func (e *guidedResponseEngine) Config(workflow Workflow) {
	// TODO -> resolve mock depedency
	flow := PokemonFlow()
	e.tree = flow
	e.node = flow.Root()
	e.setup = true
}

func (e *guidedResponseEngine) resolveMessageNode(messageID string) {

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

func (e *guidedResponseEngine) GetActionQueue() ActionQueue {
	return e.actionQueue
}

func (e *guidedResponseEngine) Name() string {
	return "guided"
}

func (e *guidedResponseEngine) HandleMessage(input models.Message, output chan<- models.Message) {

	if !e.setup {
		logger.Logger.Error("engine is not ready", zap.String("context", "guided_engine"))
		response := &input
		response.Output = "some error ocurred, please contant admin"
		output <- *response
		return
	}

	e.resolveMessageNode(input.Input)

	if e.node.Message().HasAction() {
		actionPair := ActionReplyPair{replyTo: output, action: e.node.message.Action, input: input}
		queue := e.GetActionQueue()
		queue.Put(actionPair)
	}

	options := make([]string, 0)
	e.node.TransverseInChildren(func(child *MessageNode) {
		options = append(options, child.Message().ID())
	})

	response := &input
	response.Output = e.node.Message().Content
	response.Options = options
	response.ResponseType = models.OptionResponse
	output <- *response
}

func PokemonFlow() *MessageTree {
	getPikachu := actions.NewHttpGetAction(
		"https://pokeapi.co/api/v2/pokemon/pikachu",
		"",
		&actions.TagAcess{Key: "abilities[1].ability.name"})

	getCharizard := actions.NewHttpGetAction(
		"https://pokeapi.co/api/v2/pokemon/charizard",
		"",
		&actions.TagAcess{Key: "abilities[1].ability.name"})

	tree := &MessageTree{}
	tree.Insert(
		&MessageNode{
			message: Message{
				parent:  "",
				id:      "parent",
				Content: "A habilidade de qual pokemon voce gostaria de saber?",
			},
		},
	).Insert(
		&MessageNode{
			message: Message{
				parent:  "parent",
				id:      "pikachu",
				Content: "Lets go, e a habilidade do pokemon mais querido do Ashe é...",
				Action:  getPikachu,
			},
		},
	).Insert(
		&MessageNode{
			message: Message{
				parent:  "parent",
				id:      "charizard",
				Content: "A habilidade do melhor de todos é...",
				Action:  getCharizard,
			},
		},
	)

	return tree
}
