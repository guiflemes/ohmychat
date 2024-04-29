package adapters

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"gopkg.in/yaml.v3"

	"oh-my-chat/src/core"
	"oh-my-chat/src/logger"
	"oh-my-chat/src/models"
	"oh-my-chat/src/schemas"
)

type LocalFileRepository struct {
	rootDir string
}

func NewLoadFileRepository() *LocalFileRepository {
	return &LocalFileRepository{
		rootDir: "src/examples/guided_engine/",
	}
}

func (r *LocalFileRepository) GetMessageTree(workflowID string) (*core.MessageTree, error) {

	file := workflowID + ".yml"

	log := logger.Logger.With(
		zap.String("context", "local_file_repository"),
		zap.String("file", file),
	)

	var pathTarget string

	filepath.Walk(r.rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Error("Error accessing path", zap.Error(err))
			return err
		}

		if info.Name() == file {
			pathTarget = path
			return filepath.SkipDir
		}
		return nil
	})

	if pathTarget == "" {
		log.Error("file not found")
		return nil, fmt.Errorf("file %s not found", file)
	}

	b, err := os.ReadFile(pathTarget)

	if err != nil {
		log.Error("error reading file", zap.Error(err))
		return nil, err
	}

	var model models.WorkflowGuidedModel

	if err = yaml.Unmarshal(b, &model); err != nil {
		log.Error("error Unmarshalling file into IntentModel", zap.Error(err))
		return nil, err
	}

	tree, err := schemas.ParseGuidedSchemas(model.Intents)

	if err != nil {
		log.Error("error parsing schemas", zap.Error(err))
		return nil, err
	}

	return tree, nil
}
