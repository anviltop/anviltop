package redis

import (
	"encoding/json"
	"fmt"
	"github.com/anviltop/anviltop/util"
	"github.com/garyburd/redigo/redis"
)

var (
	isInit = false
)

const (
	RUNNING_JOBS_KEY = "RUNNING_JOBS"
)

var redisConn redis.Conn

func initialize() {
	if isInit {
		return
	}

	// Setup connection - use environment variables or other configuration
	var err error
	redisConn, err = redis.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		panic("Redis dial error: " + err.Error())
	}
	isInit = true
}

func Set(key string, object interface{}) error {
	initialize()

	raw, err := json.Marshal(object)
	if err != nil {
		return err
	}

	fmt.Println("Set ", key, ". Value: ", string(raw))

	redisRes, err := redisConn.Do("SET", key, raw)
	if err != nil {
		return err
	}
	if redisRes != nil {
		fmt.Println("Redis set Result: ", util.VarDump(redisRes))
	}
	return err
}

func Get(key string, object interface{}) error {
	initialize()
	fmt.Println("Get ", key)

	value, err := redis.Bytes(redisConn.Do("GET", key))
	if err == nil {
		err = json.Unmarshal(value, object)
	}

	return err
}

func SAdd(key string, valueString string) error {
	initialize()

	_, err := redisConn.Do("SADD", key, valueString)
	if err != nil {
		panic(err)
	}
	return err
}

func SRem(key string, valueString string) error {
	initialize()

	_, err := redisConn.Do("SREM", key, valueString)
	if err != nil {
		panic(err)
	}
	return err
}

func Incr(key string, object interface) error {
	initialize()

	valueString, err := redisConn.Do("INCR", key)
	return err
}
