package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	usermodel "github.com/meneketehe/hehe/app/model/user"
	userrequest "github.com/meneketehe/hehe/app/request/user"
	userresponse "github.com/meneketehe/hehe/app/response/user"
	userservice "github.com/meneketehe/hehe/app/service/user"
)

func (h *Handler) GetMeAsUser(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	account, err := h.userAccountService.GetUserAccountByID(userID)
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
		"data":    userresponse.GetMeResponse(account.ID, account.Name, account.Email),
	})
}

func (h *Handler) RegisterUser(c *gin.Context) {
	var req userrequest.RegisterRequest
	if ok := bindData(c, &req); !ok {
		return
	}
	req.Sanitize()

	existingAccount, err := h.userAccountService.GetUserAccountByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}
	if existingAccount != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Errorf("an account with that email already exists").Error(),
			"data":    nil,
		})
		return
	}

	input := &usermodel.UserAccount{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	account, err := h.userAccountService.Register(input)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	user := usermodel.UserFromAccount(account)
	_, err = h.userService.CreateUser(user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	session := sessions.Default(c)
	session.Set("userId", account.ID)
	if err := session.Save(); err != nil {
		log.Printf("Error setting the session: %v\n", err.Error())
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "register success",
		"data":    nil,
	})
}

func (h *Handler) LoginUser(c *gin.Context) {
	var req userrequest.LoginRequest
	if ok := bindData(c, &req); !ok {
		return
	}
	req.Sanitize()

	account, err := h.userAccountService.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	session := sessions.Default(c)
	session.Set("userId", account.ID)
	if err := session.Save(); err != nil {
		log.Printf("Error setting the session: %v\n", err.Error())
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "login success",
		"data":    userresponse.GetMeResponse(account.ID, account.Name, account.Email),
	})
}

func (h *Handler) LogoutUser(c *gin.Context) {
	c.Set("userID", nil)

	session := sessions.Default(c)
	session.Set("userId", "")
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

func (h *Handler) EditUserProfile(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	var req userrequest.EditProfileRequest
	if ok := bindData(c, &req); !ok {
		return
	}
	req.Sanitize()

	account, err := h.userAccountService.GetUserAccountByID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	if !userservice.CheckPasswordHash(req.OldPassword, account.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": fmt.Errorf("old password is incorrect").Error(),
			"data":    nil,
		})
		return
	}

	account.Name = req.Name
	account.Password = req.NewPassword
	account, err = h.userAccountService.EditProfile(account)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	user.Name = req.Name
	user, err = h.userService.UpdateUser(user)
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
		"message": "edit profile success",
		"data":    nil,
	})
}
