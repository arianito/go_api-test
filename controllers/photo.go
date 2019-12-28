package ct

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/xeuus/gt/pkg/db"
	"github.com/xeuus/gt/pkg/jwt"
	"github.com/xeuus/gt/pkg/rds"
	"github.com/xeuus/instagram/api"
	"github.com/xeuus/instagram/dao"
	"image"
	_ "image/jpeg"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

const STATIC_PATH = "./static"

type Photo struct {
	RouterGroup *gin.RouterGroup
	DB          db.Database
	NAME        string
	JWT         jwt.Authenticator
	REDIS       rds.Redis
	API_ADDR    string
}

func (photo Photo) Create() {
	r := photo.RouterGroup.Group("/photo")
	helper := RedisAuthHelper{
		Prefix: photo.NAME,
		JWT:    photo.JWT,
		REDIS:  photo.REDIS,
	}

	r.GET("/image/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")

		p := dao.PhotoDAO{
			DB: photo.DB.Get(),
		}
		p.FetchByID(id)

		ctx.Header("Content-Type", p.Mime)
		ctx.File(p.URL)
	})

	r.POST("/upload", helper.Middleware, func(ctx *gin.Context) {
		userID := ctx.GetInt64("userID")

		file, err := ctx.FormFile("file0")
		if err != nil {
			panic(err)
		}

		if file.Size > 1024*1024*2 {
			panic("Entity too large")
		}
		f, err := file.Open()
		if err != nil {
			panic(err)
		}
		im, _, err := image.DecodeConfig(f)
		if err != nil {
			panic(err)
		}
		f.Close()

		ratio := float32(im.Width) / float32(im.Height)

		f, err = file.Open()
		p := dao.PhotoDAO{
			DB:     photo.DB.Get(),
			Size:   file.Size,
			GUID:   uuid.New().String(),
			URL:    "",
			InUse:  false,
			Mime:   GetFileContentType(f),
			Ratio:  ratio,
			UserID: userID,
		}
		f.Close()

		path := fmt.Sprintf("%s/%v", STATIC_PATH, userID)
		p.URL = path + "/" + p.GUID
		f, err = file.Open()

		if err != nil {
			panic(err)
		}
		data, err := ioutil.ReadAll(f)
		if err != nil {
			panic(err)
		}
		os.MkdirAll(path, 0700)
		err = ioutil.WriteFile(p.URL, data, os.ModePerm)
		if err != nil {
			panic(err)
		}
		p.DeleteUnused()
		p.Save()

		api.Success(gin.H{
			"src":  photo.API_ADDR + "/photo/image/" + p.GUID,
			"id":   p.GUID,
		}, ctx)

	})

}
func GetFileContentType(out multipart.File) string {
	buffer := make([]byte, 512)
	_, err := out.Read(buffer)
	if err != nil {
		panic(err)
	}
	contentType := http.DetectContentType(buffer)
	return contentType
}
