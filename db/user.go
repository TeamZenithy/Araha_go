package db

import (
	"context"

	"github.com/go-redis/redis/v8"

	"github.com/TeamZenithy/Araha/utils"
)

func FindGuildLocale(id string) (string, error) {
	get := utils.RDB.Get(context.Background(), "guild:locale:"+id)
	if err := get.Err(); err != nil {
		if err == redis.Nil {
			return "en", nil
		}
		return "", err
	}
	return get.Val(), nil
}

func SetGuildLocale(id, value string) error {
	return utils.RDB.Set(context.Background(), "users:locale:"+id, value, 0).Err()
}

func FindUserLocale(id string) (string, error) {
	get := utils.RDB.Get(context.Background(), "users:locale:"+id)
	if err := get.Err(); err != nil {
		if err == redis.Nil {
			return "", nil
		}
		return "", err
	}
	return get.Val(), nil
}

func SetUserLocale(id, value string) error {
	return utils.RDB.Set(context.Background(), "users:locale:"+id, value, 0).Err()
}
