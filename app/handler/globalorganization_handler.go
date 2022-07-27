package handler

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/meneketehe/hehe/app/model"
	request "github.com/meneketehe/hehe/app/request/organization"
	response "github.com/meneketehe/hehe/app/response/organization"
)

func (h *Handler) GetAllOrganizations(c *gin.Context) {
	filters := map[string]string{}
	if role, exists := c.GetQuery("role"); exists {
		filters["role"] = role
	}

	orgs, err := h.globalOrganizationService.GetAllOrganizations(filters)
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
		"data":    response.GlobalOrganizationsResponse(orgs),
	})
}

func (h *Handler) GetOrganization(c *gin.Context) {
	orgID := c.Param("ID")

	org, err := h.globalOrganizationService.GetOrganization(orgID)
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
		"data":    response.GlobalOrganizationResponse(org),
	})
}

func (h *Handler) CreateOrganization(c *gin.Context) {
	var req request.CreateOrganizationRequest
	if ok := bindData(c, &req); !ok {
		return
	}

	exists, err := h.globalOrganizationService.CheckMSPIDExists(req.MSPID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Errorf("organization with %s MSP ID already exists", req.MSPID).Error(),
			"data":    nil,
		})
		return
	}

	exists, err = h.globalOrganizationService.CheckDomainExists(req.Domain)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	if exists {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Errorf("organization with %s domain already exists", req.Domain).Error(),
			"data":    nil,
		})
		return
	}

	if req.Role == "manufacturer" {
		exists, err = h.globalOrganizationService.CheckCodeExists(req.Code)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": err.Error(),
				"data":    nil,
			})
			return
		}
		if exists {
			c.JSON(http.StatusBadRequest, gin.H{
				"success": false,
				"message": fmt.Errorf("organization with %s code already exists", req.Code).Error(),
				"data":    nil,
			})
			return
		}
	}

	input := &model.GlobalOrganization{
		Role:   req.Role,
		Name:   req.Name,
		Code:   req.Code,
		MSPID:  req.MSPID,
		Domain: req.Domain,
		Location: model.Location{
			Province:   req.Province,
			City:       req.City,
			District:   req.District,
			PostalCode: req.PostalCode,
			Address:    req.Address,
			Coordinate: model.Coordinate{
				Latitude:  req.Latitude,
				Longitude: req.Longitude,
			},
		},
		ContactInfo: model.ContactInfo{
			Phone: req.Phone,
			Email: req.Email,
		},
	}
	org, err := h.globalOrganizationService.CreateOrganization(input)
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
		"data":    response.GlobalOrganizationResponse(org),
	})
}
