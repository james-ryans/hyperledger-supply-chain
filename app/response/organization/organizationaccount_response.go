package response

type accountResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Email          string `json:"email"`
	Role           string `json:"role"`
}

func GetMeResponse(id, orgId, email, role string) *accountResponse {
	return &accountResponse{
		ID:             id,
		OrganizationID: orgId,
		Email:          email,
		Role:           role,
	}
}
