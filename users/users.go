package users

import (
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type user struct {
	Username     string `gorm:"not null;unique_index"`
	Password     string `gorm:"-"`
	PasswordHash string `gorm:"not null"`
}

func GenerateFromPassword() {

}

func CompareHashAndPassword() {

}
