package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/meneketehe/hehe/app/model"
	response "github.com/meneketehe/hehe/app/response/organization"
)

func (h *Handler) GetMeAsOrganization(c *gin.Context) {
	org, err := h.organizationService.GetMe()
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
		"data":    response.GetMeResponse(org.ID, org.Name, org.Type),
	})
}

func (h *Handler) InitOrganization(c *gin.Context) {
	orgID := os.Getenv("ORG_ID")
	orgName := os.Getenv("ORG_NAME")
	orgType := os.Getenv("ORG_TYPE")
	channelID := os.Getenv("FABRIC_CHANNEL_NAME")

	organization := model.Organization{
		ID:   orgID,
		Type: orgType,
		Name: orgName,
		Location: model.Location{
			Province:   "North Sumatra",
			City:       "Medan",
			District:   "Medan Kota",
			PostalCode: "20212",
			Address:    "Jl. Thamrin",
			Coordinate: model.Coordinate{
				Latitude:  123.1,
				Longitude: 321.1,
			},
		},
		ContactInfo: model.ContactInfo{
			Phone: "081234567890",
			Email: fmt.Sprintf("%s0@hehe.com", orgType),
		},
	}

	var err error
	switch orgType {
	case "supplier":
		_, err = h.supplierService.CreateSupplier(channelID, &model.Supplier{
			Organization: organization,
		})
	case "producer":
		_, err = h.producerService.CreateProducer(channelID, &model.Producer{
			Organization: organization,
		})
	case "manufacturer":
		_, err = h.manufacturerService.CreateManufacturer(channelID, &model.Manufacturer{
			Organization: organization,
		})
	case "distributor":
		_, err = h.distributorService.CreateDistributor(channelID, &model.Distributor{
			Vendor: model.Vendor{Organization: organization},
		})
	case "retailer":
		_, err = h.retailerService.CreateRetailer(channelID, &model.Retailer{
			Vendor: model.Vendor{Organization: organization},
		})
	}
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
		"message": fmt.Sprintf("initialize organization with %s, %s, and %s data success", orgType, orgName, orgID),
		"data":    nil,
	})
}
