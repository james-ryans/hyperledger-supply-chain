package response

type accountResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func GetMeResponse(id, name, orgType string) *accountResponse {
	return &accountResponse{
		ID:   id,
		Name: name,
		Type: orgType,
	}
}
