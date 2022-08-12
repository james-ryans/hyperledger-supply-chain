package response

import (
	"time"

	"github.com/meneketehe/hehe/app/model"
)

type getMeResponse struct {
	ID             string `json:"id"`
	OrganizationID string `json:"organization_id"`
	Code           string `json:"code"`
	Type           string `json:"type"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	Role           string `json:"role"`
}

type orgUserResponse struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	RegisteredAt time.Time `json:"registered_at"`
}

func GetMeResponse(acc *model.OrganizationAccount) *getMeResponse {
	return &getMeResponse{
		ID:             acc.ID,
		OrganizationID: acc.OrganizationID,
		Code:           acc.Code,
		Type:           acc.Type,
		Email:          acc.Email,
		Name:           acc.Name,
		Role:           acc.Role,
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
