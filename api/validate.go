package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	vstruct "github.com/xeuus/vstruct/pkg"
)

func Validate(obj interface{}, ctx ...*gin.Context) {
	vld := vstruct.NewValidator(obj)
	if len(ctx) > 0 {
		vld.BindFunc(ctx[0].Bind)
	}
	if vld.Validate(); vld.GetError() != nil {
		panic(Exception{
			Status:     http.StatusBadRequest,
			Message:    "Validation failed.",
			Errors:     []interface{}{vld.GetError().Error()},
			Validation: vld.GetMessages(),
		})
	}
}
