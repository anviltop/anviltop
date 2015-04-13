package models

import (
	"fmt"
	"github.com/anviltop/anviltop/util/redis"
)

type Model interface {
	RedisKeyPrefix() string
	Id() string
}

func redisKey(object Model, id string) string {
	result := object.RedisKeyPrefix() + "_" + id
	return result
}

func Set(object Model) error {
	redisKey := redisKey(object, object.Id())
	err := redis.Set(redisKey, object)
	return err
}

func Get(object Model, objectKey string) error {
	redisKey := redisKey(object, objectKey)
	err := redis.Get(redisKey, object)
	return err
}

func NextId(modelType string) string {
	fmt.Println("in NextId")
	var result string
	redis.Incr(modelType+"_NextId", result)
	return result
}
