package dao

import (
	gql "github.com/xeuus/gql/pkg"
)

const POSTS_TABLE = "posts"

type PostDAO struct {
	DB     interface{}
	ID     int64   `gql:"id"`
	GUID   string  `gql:"guid"`
	UserID int64   `gql:"user_id"`
	PhotoID int64   `gql:"photo_id"`
	Title    string  `gql:"title"`
}

func (post *PostDAO) Save() {
	a := gql.Create(POSTS_TABLE).
		Bind(&post).
		Use(post.DB).
		Run()
	if a.GetError() != nil {
		panic(a.GetError())
	}
}