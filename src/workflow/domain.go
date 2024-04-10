package workflow

import (
	"github.com/google/uuid"
)

type WorkFlowCatalog struct {
	Id     uuid.UUID
	Name   string
	UserId uuid.UUID
}

func (w *WorkFlowCatalog) GetWorkFlows() []WorkFlow {
	return []WorkFlow{}
}

type WorkFlow struct {
	Id         uuid.UUID
	Name       string
	CatalogId  uuid.UUID
	Properties Properties
}
