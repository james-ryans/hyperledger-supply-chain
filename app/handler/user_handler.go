package handler

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	usermodel "github.com/meneketehe/hehe/app/model/user"
	userresponse "github.com/meneketehe/hehe/app/response/user"
)

func (h *Handler) GetMeAsUser(c *gin.Context) {
	userID := c.MustGet("userID").(string)

	user, err := h.userService.GetUserByID(userID)
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
		"data":    userresponse.GetMeResponse(user.ID, user.Name, user.Email),
	})
}

func (h *Handler) InitUser(c *gin.Context) {
	userID := os.Getenv("USER_ID")
	userName := os.Getenv("USER_NAME")
	userEmail := os.Getenv("USER_EMAIL")
	userPassword := os.Getenv("USER_PASSWORD")

	hashedPassword, err := h.userService.HashPassword(userPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return
	}

	user := usermodel.User{
		ID:       userID,
		Name:     userName,
		Email:    userEmail,
		Password: hashedPassword,
	}

	_, err = h.userService.CreateUser(&user)
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
		"message": fmt.Sprintf("iniatialize user %s success", userName),
		"data":    nil,
	})
}
