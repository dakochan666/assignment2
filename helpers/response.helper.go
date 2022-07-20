package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func BadRequestResponse(ctx *gin.Context, payload interface{}) {
	WriteJsonResponse(ctx, http.StatusBadRequest, gin.H{
		"success": false,
		"error":   payload,
	})
}

func NotFoundResponse(ctx *gin.Context, payload interface{}) {
	WriteJsonResponse(ctx, http.StatusNotFound, gin.H{
		"success": false,
		"error":   payload,
	})
}

func WriteJsonResponse(ctx *gin.Context, status int, payload interface{}) {
	ctx.JSON(status, payload)
}
