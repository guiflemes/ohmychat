package models

type ModelType string

const (
	TypeHttpGetModel     ModelType = "http_get"
	TypeGoogleSheetModel ModelType = "google_sheet"
)
