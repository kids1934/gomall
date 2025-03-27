package dal

import (
	"github.com/kids1934/gomall/app/order/biz/dal/mysql"
	"github.com/kids1934/gomall/app/order/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
