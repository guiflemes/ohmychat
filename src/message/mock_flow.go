package message

import (
	"oh-my-chat/src/actions"
)

func NotionFlow() *MessageTree {
	n := actions.NewNotionActions().PedencyGetter
	tree := &MessageTree{}
	tree.Insert(&MessageNode{message: Message{parent: "", id: "coco", Content: "O que voce gostaria de saber?"}}).
		Insert(&MessageNode{message: Message{parent: "coco", id: "tarefas", Content: "Tarefas, escolha as opções"}}).
		Insert(&MessageNode{message: Message{parent: "coco", id: "assinaturas", Content: "Assinaturas, esolhas as opções"}}).
		Insert(&MessageNode{message: Message{parent: "coco", id: "marvin", Content: "Marvin, escolha o role"}}).
		Insert(&MessageNode{message: Message{parent: "tarefas", id: "pendencias", Content: "ok, verificando", Action: n}}).
		Insert(&MessageNode{message: Message{parent: "tarefas", id: "pagas"}}).
		Insert(&MessageNode{message: Message{parent: "marvin", id: "coco"}})

	return tree
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
