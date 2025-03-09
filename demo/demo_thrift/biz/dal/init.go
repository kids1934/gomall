package dal

import (
	"github.com/kids1934/gomall/demo/demo_thrift/biz/dal/mysql"
	"github.com/kids1934/gomall/demo/demo_thrift/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
