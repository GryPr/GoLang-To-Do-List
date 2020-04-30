package users

import (
	"github.com/jinzhu/gorm"
)

var db *gorm.DB

// InitDB starts the database for the package
func InitDB(psql string) {
	var err error
	db, err = gorm.Open("postgres", psql)
	if err != nil {
		panic(err)
	}
}
