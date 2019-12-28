package dao

import (
	"fmt"
	gql "github.com/xeuus/gql/pkg"
	"time"
)

type FeedDAO struct {
	DB interface{} `json:"-"`

	PostGUID    string    `gql:"post_guid" json:"id"`
	Username    string    `gql:"username" json:"username"`
	Title       string    `gql:"title" json:"title"`
	PhotoGUID   string    `gql:"photo_guid" json:"-"`
	PhotoRatio  float32   `gql:"photo_ratio" json:"ratio"`
	CreatedDate time.Time `gql:"created_date" json:"createdDate"`
	Src         string    `json:"src"`
}

func (feed *FeedDAO) FetchList() []*FeedDAO {
	var list []*FeedDAO
	if a := gql.Read("posts ps").
		Model(&feed).
		Columns(
			"ps.title title",
			"ps.guid post_guid",
			"us.username username",
			"ph.guid photo_guid",
			"ph.photo_ratio photo_ratio",
			"ps.created_date created_date",
		).
		Join("users us", "ps.user_id = us.id").
		Join("photos ph", "ps.photo_id = ph.id").
		OrderBy("-ps.created_date").
		Use(feed.DB).
		Scan(&list); a.GetError() != nil {
		fmt.Println(a.Query())
		panic(a.GetError())
	}
	return list
}
