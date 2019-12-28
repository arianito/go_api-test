package ct

import (
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/google/uuid"
	"github.com/xeuus/gt/pkg/db"
	"github.com/xeuus/gt/pkg/hash"
	"github.com/xeuus/gt/pkg/jwt"
	"github.com/xeuus/gt/pkg/rds"
	"github.com/xeuus/instagram/api"
	"github.com/xeuus/instagram/dao"
	"net/http"
	"strconv"
	"time"
)

type Auth struct {
	RouterGroup *gin.RouterGroup
	DB          db.Database
	JWT         jwt.Authenticator
	HASH        hash.Hash
	REDIS       rds.Redis
	NAME        string
}

func (auth Auth) Create() {
	timeout := 120
	r := auth.RouterGroup.Group("/auth")
	r.POST("/login/otp", func(ctx *gin.Context) {
		type LoginRequest struct {
			Username string `json:"username" v:"required mobile_iran"`
			OTP      string `json:"otp" v:"required min(4)"`
		}
		obj := new(LoginRequest)
		api.Validate(&obj, ctx)

		clientName := ctx.GetHeader("Client-Name")
		gateway := ctx.GetHeader("App-Gateway")
		code := ""
		helper := RedisAuthHelper{
			Prefix:  auth.NAME,
			Timeout: timeout,
			JWT:     auth.JWT,
		}
		err := auth.REDIS.Action(func(c redis.Conn) error {
			var err error
			helper.Conn = c
			code, err = helper.Get(obj.Username)
			return err
		})
		if err != nil {
			panic("Wrong code")
		}
		if code == "" || code != obj.OTP {
			panic("Wrong code")
		}
		user := dao.UserDAO{
			DB: auth.DB.Get(),
		}
		user.Fetch(obj.Username)
		if !user.Active {
			panic("User inactive")
		}
		token := ""
		err = auth.REDIS.Action(func(c redis.Conn) error {
			var err error
			helper.Conn = c
			token, err = helper.Create(strconv.FormatInt(user.ID, 10), clientName, gateway, time.Hour*24, 0)
			if err != nil {
				return err
			}
			return helper.Remove(obj.Username)
		})
		if err != nil {
			panic(err)
		}
		api.Success(gin.H{
			"token": token,
		}, ctx)
	})
	r.POST("/login/password", func(ctx *gin.Context) {
		type LoginRequest struct {
			Username string `json:"username" v:"required"`
			Password string `json:"password" v:"required min(5)"`
		}
		obj := new(LoginRequest)
		api.Validate(&obj, ctx)

		clientName := ctx.GetHeader("Client-Name")
		gateway := ctx.GetHeader("App-Gateway")
		helper := RedisAuthHelper{
			Prefix:  auth.NAME,
			Timeout: timeout,
			JWT:     auth.JWT,
		}

		user := dao.UserDAO{
			DB: auth.DB.Get(),
		}
		user.Fetch(obj.Username)
		if !user.Active || !auth.HASH.CheckHash(user.Password, obj.Password) {
			panic("Invalid credentials")
		}
		token := ""
		err := auth.REDIS.Action(func(c redis.Conn) error {
			var err error
			helper.Conn = c
			token, err = helper.Create(strconv.FormatInt(user.ID, 10), clientName, gateway, time.Hour*24, 0)
			if err != nil {
				return err
			}
			return helper.Remove(obj.Username)
		})
		if err != nil {
			panic(err)
		}
		api.Success(gin.H{
			"token": token,
		}, ctx)
	})
	r.GET("/otp/:mobileNumber", func(ctx *gin.Context) {
		type MobileNumberRequest struct {
			MobileNumber string `json:"mobileNumber" v:"required mobile_iran"`
		}
		obj := new(MobileNumberRequest)
		obj.MobileNumber = ctx.Param("mobileNumber")
		api.Validate(&obj)
		helper := RedisAuthHelper{
			Prefix:  auth.NAME,
			Timeout: timeout,
			JWT:     auth.JWT,
			KeyGen:  auth.HASH.GenerateFixed,
		}
		remaining := 0
		if err := auth.REDIS.Action(func(conn redis.Conn) error {
			helper.Conn = conn
			var err error
			remaining, err = helper.Request(obj.MobileNumber, func(s string) error {
				return nil
			})
			if err != nil {
				return err
			}
			return nil
		}); err != nil {
			panic(err)
		}
		api.Success(gin.H{"resendAt": remaining}, ctx)
	})
	r.POST("/check", func(ctx *gin.Context) {
		type CheckRequest struct {
			Username     string `json:"username" v:"required username"`
			MobileNumber string `json:"mobileNumber" v:"required mobile_iran"`
			Password     string `json:"password" v:"required min(5)"`
			Name         string `json:"name" v:"required"`
		}
		obj := new(CheckRequest)
		api.Validate(&obj, ctx)

		user := dao.UserDAO{
			DB:           auth.DB.Get(),
			Username:     obj.Username,
			MobileNumber: obj.MobileNumber,
		}
		if user.Exists() {
			panic(api.Exception{
				Status:  http.StatusBadRequest,
				Message: "User already exists.",
				Validation: map[string]string{
					"username":     "Username already exist.",
					"mobileNumber": "Mobile Number already owned.",
				},
			})
		}
		api.Success(gin.H{}, ctx)
	})
	r.POST("/register", func(ctx *gin.Context) {
		type RegisterRequest struct {
			Username     string `json:"username" v:"required username"`
			MobileNumber string `json:"mobileNumber" v:"required mobile_iran"`
			Password     string `json:"password" v:"required min(5)"`
			Name         string `json:"name" v:"required"`
			OTP          string `json:"otp" v:"required min(4)"`
		}
		obj := new(RegisterRequest)
		api.Validate(&obj, ctx)

		clientName := ctx.GetHeader("Client-Name")
		gateway := ctx.GetHeader("App-Gateway")

		code := ""
		helper := RedisAuthHelper{
			Prefix:  auth.NAME,
			Timeout: timeout,
			JWT:     auth.JWT,
		}
		err := auth.REDIS.Action(func(c redis.Conn) error {
			var err error
			helper.Conn = c
			code, err = helper.Get(obj.MobileNumber)
			return err
		})
		if err != nil {
			panic("Wrong code")
		}
		if code == "" || code != obj.OTP {
			panic("Wrong code")
		}
		user := dao.UserDAO{
			DB:           auth.DB.Get(),
			GUID:         uuid.New().String(),
			Username:     obj.Username,
			MobileNumber: obj.MobileNumber,
			Name:         obj.Name,
			Password:     auth.HASH.Hash(obj.Password),
			Active:       true,
			Superuser:    false,
			Bio:          "",
		}
		user.Save()
		token := ""
		err = auth.REDIS.Action(func(c redis.Conn) error {
			var err error
			helper.Conn = c
			token, err = helper.Create(strconv.FormatInt(user.ID, 10), clientName, gateway, time.Hour*24, 0)
			if err != nil {
				return err
			}
			return helper.Remove(obj.Username)
		})
		if err != nil {
			panic(err)
		}
		api.Success(gin.H{
			"token": token,
		}, ctx)
	})

}
