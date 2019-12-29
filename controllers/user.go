package ct

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xeuus/gt/pkg/db"
	"github.com/xeuus/gt/pkg/jwt"
	"github.com/xeuus/gt/pkg/rds"
	"github.com/xeuus/instagram/api"
	"github.com/xeuus/instagram/dao"
)

type User struct {
	RouterGroup *gin.RouterGroup
	DB          db.Database
	NAME        string
	JWT         jwt.Authenticator
	REDIS       rds.Redis
}

func (user User) Create() {
	r := user.RouterGroup.Group("/user")
	helper := RedisAuthHelper{
		Prefix: user.NAME,
		JWT:    user.JWT,
		REDIS:  user.REDIS,
	}
	r.Use(helper.Middleware)

	r.GET("/info", func(ctx *gin.Context) {
		userID := ctx.GetInt64("userID")
		u := dao.UserDAO{
			DB: user.DB.Get(),
		}
		u.FetchByID(userID)
		api.Success(gin.H{
			"id":           userID,
			"name":         u.Name,
			"username":     u.Username,
			"mobileNumber": u.MobileNumber,
		}, ctx)
	})

	r.POST("/post", func(ctx *gin.Context) {
		type PostRequest struct {
			Title   string `json:"title"`
			PhotoId string `json:"photoId" v:"required"`
		}
		obj := new(PostRequest)
		api.Validate(&obj, ctx)
		userID := ctx.GetInt64("userID")

		err := user.DB.Transaction(func(tx *sql.Tx) error {

			p := dao.PhotoDAO{
				DB: tx,
			}
			p.FetchByID(obj.PhotoId)
			p.InUse = true
			p.Update()

			po := dao.PostDAO{
				DB:      tx,
				UserID:  userID,
				GUID:    uuid.New().String(),
				PhotoID: p.ID,
				Title:   obj.Title,
			}
			po.Save()

			return nil
		})
		if err != nil {
			panic(err)
		}
		api.Success(gin.H{}, ctx)

	})

}
