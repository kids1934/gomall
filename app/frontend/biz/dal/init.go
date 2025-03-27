package dal

import (
	"github.com/kids1934/gomall/app/frontend/biz/dal/mysql"
)

func Init() {
	// redis.Init()
	mysql.Init()
}
