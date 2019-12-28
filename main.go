package main

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	mgo "github.com/xeuus/amigo/pkg"
	"github.com/xeuus/gt/pkg/db"
	"github.com/xeuus/instagram/api"
	"github.com/xeuus/instagram/controllers"
	"github.com/xeuus/vstruct/pkg"
	"log"
)

func main() {
	vstruct.LoadBuiltin()
	DB := db.NewClient("mysql", DB_QUERY)
	DB.Connect()
	defer DB.Close()
	mgo.SetTable("migrations")
	mgo.Migrate("./migrations", "up", "", DB.Get())
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(api.PanicMiddleware)
	router := r.Group(API_PREFIX)
	ct.Feed{
		RouterGroup: router,
		JWT:         JWT,
		DB:          DB,
		NAME:        NAME,
		REDIS:       REDIS,
		API_ADDR:       API_ADDR,
	}.Create()
	ct.Auth{
		RouterGroup: router,
		JWT:         JWT,
		DB:          DB,
		HASH:        HASH,
		NAME:        NAME,
		REDIS:       REDIS,
	}.Create()
	ct.User{
		RouterGroup: router,
		JWT:         JWT,
		DB:          DB,
		NAME:        NAME,
		REDIS:       REDIS,
	}.Create()
	ct.Photo{
		RouterGroup: router,
		JWT:         JWT,
		DB:          DB,
		NAME:        NAME,
		REDIS:       REDIS,
		API_ADDR:       API_ADDR,
	}.Create()
	log.Println("[Http] Listening on", PORT)
	if err := r.Run(PORT); err != nil {
		log.Fatal(err)
	}
}
