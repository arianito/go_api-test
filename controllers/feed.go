package ct

import (
	"github.com/gin-gonic/gin"
	"github.com/xeuus/gt/pkg/db"
	"github.com/xeuus/gt/pkg/jwt"
	"github.com/xeuus/gt/pkg/rds"
	"github.com/xeuus/instagram/api"
	"github.com/xeuus/instagram/dao"
)

type Feed struct {
	RouterGroup *gin.RouterGroup
	DB          db.Database
	NAME        string
	JWT         jwt.Authenticator
	REDIS       rds.Redis
	API_ADDR    string
}

func (feed Feed) Create() {

	r := feed.RouterGroup.Group("/feed")

	helper := RedisAuthHelper{
		Prefix: feed.NAME,
		JWT:    feed.JWT,
		REDIS:  feed.REDIS,
	}
	r.Use(helper.Middleware)

	r.GET("/list", func(ctx *gin.Context) {
		fd := dao.FeedDAO{
			DB: feed.DB.Get(),
		}
		lst := fd.FetchList()
		for _, a := range lst {
			a.Src = feed.API_ADDR + "/photo/image/" + a.PhotoGUID
		}
		api.Success(lst, ctx)
	})

}
