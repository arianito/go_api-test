package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func PanicMiddleware(ctx *gin.Context) {
	defer func() {
		if p := recover(); p != nil {
			stack := "" // string(debug.Stack())
			switch p.(type) {
			case string:
				ctx.JSON(http.StatusInternalServerError, Exception{
					Status:     http.StatusInternalServerError,
					Message:    p.(string),
					Errors:     []interface{}{},
					Validation: map[string]string{},
					Stack:      stack,
				})
				break
			case error:
				err := p.(error)
				ctx.JSON(http.StatusInternalServerError, Exception{
					Status:     http.StatusInternalServerError,
					Message:    err.Error(),
					Errors:     []interface{}{},
					Validation: map[string]string{},
					Stack:      stack,
				})
				break
			case Exception:
				api := p.(Exception)
				api.Stack = stack
				ctx.JSON(api.Status, api)
				break
			default:
				ctx.JSON(http.StatusInternalServerError, Exception{
					Status:     http.StatusInternalServerError,
					Message:    "خطای سمت سرور",
					Errors:     []interface{}{p},
					Validation: map[string]string{},
					Stack:      stack,
				})
				break
			}
			ctx.Abort()
		}
	}()
	ctx.Next()
}