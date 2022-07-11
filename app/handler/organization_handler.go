package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meneketehe/hehe/app/model"
)

func (h *Handler) GetMeAsOrganization(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	organization, err := h.organizationService.GetMe(orgID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": nil,
		"data":    model.ToOrganizationResponse(organization),
	})
}
