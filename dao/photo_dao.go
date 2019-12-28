package dao

import (
	gql "github.com/xeuus/gql/pkg"
	"os"
)

const PHOTOS_TABLE = "photos"

type PhotoDAO struct {
	DB     interface{}
	ID     int64   `gql:"id"`
	GUID   string  `gql:"guid"`
	UserID int64   `gql:"user_id"`
	URL    string  `gql:"photo_url"`
	Ratio  float32 `gql:"photo_ratio"`
	Mime   string  `gql:"photo_mime"`
	Size   int64   `gql:"photo_size"`
	InUse  bool    `gql:"in_use"`
}
func (photo *PhotoDAO) FetchByID(id string) {
	if a := gql.Read(PHOTOS_TABLE).
		Model(&photo).
		Where("GUID", id).
		Use(photo.DB).
		Scan(&photo); !a.HasValue() {
		panic("Not found")
	}
}


func (photo *PhotoDAO) Save() {
	a := gql.Create(PHOTOS_TABLE).
		Bind(&photo).
		Use(photo.DB).
		Run()
	if a.GetError() != nil {
		panic(a.GetError())
	}
}


func (photo *PhotoDAO) Update() {
	a := gql.Update(PHOTOS_TABLE).
		Bind(&photo).
		Where("ID", photo.ID).
		Use(photo.DB).
		Run()
	if a.GetError() != nil {
		panic(a.GetError())
	}
}

func (photo *PhotoDAO) DeleteUnused() {
	var photos []*PhotoDAO
	if a := gql.Read(PHOTOS_TABLE).
		Model(&photo).
		Where("UserID", photo.UserID).
		Where("InUse", false).
		Use(photo.DB).
		Scan(&photos); a.GetError() != nil {
		panic(a.GetError())
	}
	for _, p := range photos {
		_ = os.Remove(p.URL)
	}

	if a := gql.Delete(PHOTOS_TABLE).
		Model(&photo).
		Where("UserID", photo.UserID).
		Where("InUse", false).
		Use(photo.DB).
		Run(); a.GetError() != nil {
		panic(a.GetError())
	}
}
