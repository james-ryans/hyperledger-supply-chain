package response

import (
	"time"

	"github.com/meneketehe/hehe/app/model"
)

type getMeResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Email          string `json:"email"`
	Role           string `json:"role"`
}

type orgUserResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	RegisteredAt time.Time `json:"registered_at"`
}

func GetMeResponse(id, orgId, email, role string) *getMeResponse {
	return &getMeResponse{
		ID:             id,
		OrganizationID: orgId,
		Email:          email,
		Role:           role,
	}
}

func OrgUsersResponse(accs []*model.OrganizationAccount) []*orgUserResponse {
	res := make([]*orgUserResponse, 0)
	for _, acc := range accs {
		res = append(res, OrgUserResponse(acc))
	}

	return res
}

func OrgUserResponse(acc *model.OrganizationAccount) *orgUserResponse {
	return &orgUserResponse{
		ID:           acc.ID,
		Name:         acc.Name,
		Email:        acc.Email,
		Phone:        acc.Phone,
		RegisteredAt: acc.RegisteredAt,
	}
}
