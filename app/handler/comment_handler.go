package handler

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	usermodel "github.com/meneketehe/hehe/app/model/user"
	userrequest "github.com/meneketehe/hehe/app/request/user"
	userresponse "github.com/meneketehe/hehe/app/response/user"
)

func (h *Handler) GetComments(c *gin.Context) {
	riceID := c.Param("riceID")

	comments, err := h.commentService.GetAllCommentByRiceID(riceID)
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
		"data":    userresponse.CommentsResponse(comments),
	})
}

func (h *Handler) WriteComment(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	riceID := c.Param("riceID")

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": fmt.Errorf("user not found").Error(),
			"data":    nil,
		})
		return
	}

	var req userrequest.WriteCommentRequest
	if ok := bindData(c, &req); !ok {
		return
	}

	input := &usermodel.Comment{
		RiceID:    riceID,
		UserName:  user.Name,
		Text:      req.Text,
		CommentAt: time.Now(),
	}
	_, err = h.commentService.WriteComment(input)
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
		"message": "write comment success",
		"data":    nil,
	})
}
