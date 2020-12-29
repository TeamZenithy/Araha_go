package db

import (
	"github.com/TeamZenithy/Araha/logger"
	"github.com/TeamZenithy/Araha/utils"
	"github.com/go-redis/redis/v8"
)

//utils.RedisHost, RedisPass, RedisPort

func InitRedis() {

	opt, err := redis.ParseURL("redis://" + utils.RedisHost + ":" + utils.RedisPort + "")
	if err != nil {
		logger.Panic(err.Error())
	}

	opt.Password = utils.RedisPass

	rdb := redis.NewClient(opt)

	utils.RDB = rdb
}
