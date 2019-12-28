package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Success(obj interface{}, ctx *gin.Context) {
	ctx.JSON(http.StatusOK, obj)
}
