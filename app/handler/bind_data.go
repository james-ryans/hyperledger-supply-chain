package handler

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/meneketehe/hehe/app/model"
)

type Request interface {
	Validate() error
}

func bindData(c *gin.Context, req Request) bool {
	if err := c.ShouldBind(req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
			"data":    nil,
		})
		return false
	}

	if err := req.Validate(); err != nil {
		errs := strings.Split(err.Error(), ";")
		errors := make([]model.FieldError, 0)

		for _, e := range errs {
			split := strings.Split(e, ":")
			er := model.FieldError{
				Field:   strings.TrimSpace(split[0]),
				Message: strings.TrimSpace(split[1]),
			}
			errors = append(errors, er)
		}

		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"success": false,
			"message": "Data validation error.",
			"data":    model.ToFieldErrorsResponse(&errors),
		})
		return false
	}

	return true
}
