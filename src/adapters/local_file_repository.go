package adapters

import (
	"oh-my-chat/src/core"
)

type LocalFileRepository struct{}

func (r *LocalFileRepository) GetMessageTree(workflowID string) (*core.MessageTree, error) {
	return core.PokemonFlow(), nil
}
