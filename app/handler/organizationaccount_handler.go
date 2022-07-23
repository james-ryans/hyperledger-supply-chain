package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	request "github.com/meneketehe/hehe/app/request/organization"
	response "github.com/meneketehe/hehe/app/response/organization"
	service "github.com/meneketehe/hehe/app/service/organization"
)

func (h *Handler) GetMeAsOrganization(c *gin.Context) {
	adminID := c.MustGet("adminID").(string)

	account, err := h.organizationAccountService.GetOrganizationAccountByID(adminID)
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
		"data":    response.GetMeResponse(account.ID, account.OrganizationID, account.Email, account.Role),
	})
}

func (h *Handler) LoginOrganization(c *gin.Context) {
	var req request.LoginRequest
	if ok := bindData(c, &req); !ok {
		return
	}
	req.Sanitize()

	log.Printf("%s %s\n", req.Email, req.Password)

	account, err := h.organizationAccountService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	session := sessions.Default(c)
	session.Set("adminId", account.ID)
	session.Set("orgId", account.OrganizationID)
	if err := session.Save(); err != nil {
		log.Printf("Error setting the session: %v\n", err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "login success",
		"data":    response.GetMeResponse(account.ID, account.OrganizationID, account.Email, account.Role),
	})
}

func (h *Handler) LogoutOrganization(c *gin.Context) {
	c.Set("adminID", nil)
	c.Set("orgID", nil)

	session := sessions.Default(c)
	session.Set("adminId", "")
	session.Set("orgId", "")
	session.Clear()
	session.Options(sessions.Options{Path: "/", MaxAge: -1})
	if err := session.Save(); err != nil {
		log.Printf("Error clearing session: %v\n", err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "logout success",
		"data":    nil,
	})
}

func (h *Handler) ChangePasswordOrganization(c *gin.Context) {
	adminID := c.MustGet("adminID").(string)

	var req request.ChangePasswordRequest
	if ok := bindData(c, &req); !ok {
		return
	}
	req.Sanitize()

	account, err := h.organizationAccountService.GetOrganizationAccountByID(adminID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if !service.CheckPasswordHash(req.OldPassword, account.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Errorf("old password is incorrect").Error(),
			"data":    nil,
		})
		return
	}

	account.Password = req.NewPassword
	account, err = h.organizationAccountService.ChangePassword(account)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "change password success",
		"data":    nil,
	})
}
