package handler

import (
	"ai-tutor-backend/internal/apperr"

	"github.com/gin-gonic/gin"
)

func writeError(c *gin.Context, err error) {
	appErr := apperr.From(err)

	c.JSON(appErr.Status, gin.H{
		"error_message": appErr.Error(),
	})
 

}