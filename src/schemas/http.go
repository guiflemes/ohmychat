package schemas

type HttpGetSchema struct {
	ID        string     `json:"id,omitempty"`
	Type      SchemaType `json:"type,omitempty"`
	Url       string     `json:"url,omitempty"`
	Auth      string     `json:"auth"`
	TagAccess string     `json:"tag_access"`
}

func (h HttpGetSchema) GetID() string {
	return h.ID
}
func (h HttpGetSchema) GetType() SchemaType {
	return h.Type
}
