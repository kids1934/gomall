package dal

import (
	"github.com/kids1934/gomall/app/email/biz/dal/mysql"
	"github.com/kids1934/gomall/app/email/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
