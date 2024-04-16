package models

type HttpGetProperty struct {
	ID        string       `json:"id,omitempty"`
	Type      PropertyType `json:"type,omitempty"`
	Url       string       `json:"url,omitempty"`
	Auth      string       `json:"auth"`
	TagAccess string       `json:"tag_access"`
}

func (h HttpGetProperty) GetID() string {
	return h.ID
}
func (h HttpGetProperty) GetType() PropertyType {
	return h.Type
}
