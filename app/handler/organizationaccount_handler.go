package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/meneketehe/hehe/app/model"
	"github.com/meneketehe/hehe/app/model/enum"
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
	session.Set("role", account.Role)
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
	c.Set("role", nil)

	session := sessions.Default(c)
	session.Set("adminId", "")
	session.Set("orgId", "")
	session.Set("role", "")
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

func (h *Handler) GetAllUsers(c *gin.Context) {
	accs, err := h.organizationAccountService.GetAllOrganizationUserAccounts()
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
		"data":    response.OrgUsersResponse(accs),
	})
}

func (h *Handler) GetUser(c *gin.Context) {
	ID := c.Param("ID")

	acc, err := h.organizationAccountService.GetOrganizationAccountByID(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	if acc == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"message": fmt.Errorf("user not found").Error(),
			"data":    nil,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": nil,
		"data":    response.OrgUserResponse(acc),
	})
}

func (h *Handler) CreateUser(c *gin.Context) {
	orgID := c.MustGet("orgID").(string)
	role := c.MustGet("role").(string)

	var req request.CreateUserRequest
	if ok := bindData(c, &req); !ok {
		return
	}
	req.Sanitize()

	acc, err := h.organizationAccountService.GetOrganizationAccountByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	if acc != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Errorf("an account with that email already exists").Error(),
			"data":    nil,
		})
		return
	}

	input := &model.OrganizationAccount{
		Role:           role,
		Type:           enum.OrgAccUser,
		OrganizationID: orgID,
		Name:           req.Name,
		Email:          req.Email,
		Phone:          req.Phone,
		Password:       req.Password,
		RegisteredAt:   time.Now(),
	}

	acc, err = h.organizationAccountService.CreateUser(input)
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
		"message": "create user success",
		"data":    response.OrgUserResponse(acc),
	})
}

func (h *Handler) UpdateUser(c *gin.Context) {
	var req request.UpdateUserRequest
	if ok := bindData(c, &req); !ok {
		return
	}
	req.Sanitize()

	ID := c.Param("ID")
	acc, err := h.organizationAccountService.GetOrganizationAccountByID(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	acc.Name = req.Name
	acc.Phone = req.Phone
	if req.ChangePassword {
		acc.Password = req.Password
	}

	acc, err = h.organizationAccountService.UpdateUser(acc, req.ChangePassword)
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
		"message": "update user success",
		"data":    response.OrgUserResponse(acc),
	})
}

func (h *Handler) DeleteUser(c *gin.Context) {
	ID := c.Param("ID")
	acc, err := h.organizationAccountService.GetOrganizationAccountByID(ID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	err = h.organizationAccountService.DeleteUser(acc)
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
		"message": "delete user success",
		"data":    nil,
	})
}
