package ct

import (
	"github.com/gin-gonic/gin"
	gql "github.com/xeuus/gql/pkg"
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

	r.POST("/list", func(ctx *gin.Context) {
		type ListRequest struct {
			LastTime gql.NullTime  `json:"lastTime"`
		}
		obj := new(ListRequest)
		api.Validate(&obj, ctx)

		fd := dao.FeedDAO{
			DB: feed.DB.Get(),
		}
		lst := fd.FetchList(obj.LastTime)
		for _, a := range lst {
			a.Src = feed.API_ADDR + "/photo/image/" + a.PhotoGUID
		}
		api.Success(lst, ctx)
	})

}
