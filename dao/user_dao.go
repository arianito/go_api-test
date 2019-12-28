package dao

import (
	gql "github.com/xeuus/gql/pkg"
)

const USERS_TABLE = "users"

type UserDAO struct {
	DB           interface{}
	ID           int64  `gql:"id"`
	Name         string `gql:"name"`
	GUID         string `gql:"guid"`
	Bio          string `gql:"bio"`
	Username     string `gql:"username"`
	MobileNumber string `gql:"mobile_number"`
	Password     string `gql:"password"`
	Active       bool   `gql:"active"`
	Superuser    bool   `gql:"superuser"`
}

func (user *UserDAO) Fetch(username string) {
	if a := gql.Read(USERS_TABLE).
		Model(&user).
		Where("Username", username).
		Or().
		Where("MobileNumber", username).
		Use(user.DB).
		Scan(&user);!a.HasValue() {
		panic("Not found")
	}
}

func (user *UserDAO) FetchByID(id int64) {
	if a := gql.Read(USERS_TABLE).
		Model(&user).
		Where("ID", id).
		Use(user.DB).
		Scan(&user); !a.HasValue() {
		panic("Not found")
	}
}

func (user *UserDAO) Exists() bool {
	a := gql.Read(USERS_TABLE).
		Model(&user).
		Where("Username", user.Username).
		Or().
		Where("MobileNumber", user.MobileNumber).
		Use(user.DB).
		Scan(&user)
	if !a.HasValue() {
		return false
	}
	return true
}

func (user *UserDAO) Save() {
	a := gql.Create(USERS_TABLE).
		Bind(&user).
		Use(user.DB).
		Run()
	if a.GetError() != nil {
		panic(a.GetError())
	}
}
