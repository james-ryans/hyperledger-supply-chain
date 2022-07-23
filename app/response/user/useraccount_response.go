package userresponse

type accountResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func GetMeResponse(id, name, email string) *accountResponse {
	return &accountResponse{
		ID:    id,
		Name:  name,
		Email: email,
	}
}
