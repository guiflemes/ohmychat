package models

type CollumnType int

const (
	CollumnTypeText CollumnType = 1
)

type GoogleSheetModel struct {
	SecretPath  string      `json:"secret_path"`
	WriteConfig WriteConfig `json:"write_config"`
}

type WriteConfig struct {
	CollumnName string      `json:"collumn_name"`
	CollumnType CollumnType `json:"collumn_type"`
}

func (h *GoogleSheetModel) GetType() string {
	return "http_get"
}
