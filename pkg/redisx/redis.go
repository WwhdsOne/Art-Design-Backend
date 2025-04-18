package redisx

import (
	"Art-Design-Backend/global"
	"context"
	"time"
)

func Set(id, value string, duration time.Duration) (err error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	err = global.Redis.Set(timeout, id, value, duration).Err()
	if err != nil {
		global.Logger.Error(err.Error())
		return
	}
	return
}

func Get(key string) (val string) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	val, err := global.Redis.Get(timeout, key).Result()
	if err != nil {
		if err.Error() == "redis: nil" {
			return
		}
		global.Logger.Error(err.Error())
		return
	}
	return
}

func Delete(key string) (err error) {
	timeout, cancelFunc := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelFunc()
	err = global.Redis.Del(timeout, key).Err()
	if err != nil {
		global.Logger.Error(err.Error())
		return
	}
	return
}
