package main

import (
	"github.com/xeuus/gt/pkg/env"
	"github.com/xeuus/gt/pkg/hash"
	"github.com/xeuus/gt/pkg/jwt"
	"github.com/xeuus/gt/pkg/rds"
)

var (

	NAME = "instagram"
	PORT       = "[::]:8080"
	API_ADDR       = env.String("API_ADDR", "http://localhost:8080/api")
	API_PREFIX = env.String("API_PREFIX", "/api")
	DB_QUERY   = env.String("DB_QUERY", "develop:123@tcp(localhost:3306)/instagram_db?charset=utf8mb4&parseTime=True&loc=UTC&multiStatements=true")
	JWT        = jwt.NewClient(&jwt.Client{
		PrivateKey: env.String("PRIVATE_KEY"),
		PublicKey:  env.String("PUBLIC_KEY"),
	})
	REDIS = rds.NewClient(&rds.Client{
		Addr: env.Array("REDIS_ADDR", "localhost:6379"),
	})
	HASH = hash.NewClient()
)
