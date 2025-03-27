package main

import (
	"github.com/joho/godotenv"
	"github.com/kids1934/gomall/app/user/biz/dal"
	"github.com/kids1934/gomall/app/user/biz/dal/mysql"
	"github.com/kids1934/gomall/app/user/biz/model"
)

func FlagDB() {
	err := mysql.DB.AutoMigrate(
		&model.User{},
	)
	if err != nil {
		panic(err)
	}
}

func main() {
	_ = godotenv.Load()
	dal.Init()
	FlagDB()
}
