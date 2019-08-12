//redis配置，redis相关方法
package redis

import (
	"encoding/json"
	"errors"
	// "maoguo/henan/misc/goredis"
	"maoguo/henan/misc/config"
	"maoguo/henan/misc/utils/parse"
	"maoguo/henan/vo"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/wonderivan/logger"
)

const (
	RedisURL            = "127.0.0.1:6379"
	redisMaxIdle        = 10
	redisIdleTimeoutSec = 240
	maxActive           = 1024
	// RedisPassword       = config.CONFIG[""]
	// RedisPassword = "123456"
)

func NewRedisPool() *redis.Pool {
	logger.Info("redis pwd:", config.CONFIG["RedisPassword"])
	v := &redis.Pool{
		MaxIdle:     redisMaxIdle,
		MaxActive:   maxActive,
		IdleTimeout: 1 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", RedisURL,
				redis.DialPassword(config.CONFIG["RedisPassword"]),
				redis.DialDatabase(int(0)),
				redis.DialConnectTimeout(1*time.Second),
				redis.DialReadTimeout(1*time.Second),
				redis.DialWriteTimeout(1*time.Second))
			if err != nil {
				logger.Error("redis conn failed,", err)
				return nil, err
			}
			return con, nil
		},
	}
	return v
}

var rc *redis.Pool

func RedisClientGet() *redis.Pool {
	if rc == nil {
		rc = NewRedisPool()
	}
	return rc
}

func Set(k, v string) {
	c := RedisClientGet().Get()
	defer c.Close()
	_, err := c.Do("SET", k, v)
	if err != nil {
		logger.Error("set error", err.Error())
	}
}

func SetT(k, v string) {
	c := RedisClientGet().Get()
	defer c.Close()
	_, err := c.Do("SET", []byte(k), []byte(v))
	if err != nil {
		logger.Error("set error", err.Error())
	}
}

func HGet(mapname, k string) string {
	c := RedisClientGet().Get()
	defer c.Close()
	r, err := redis.String(c.Do("HGet", mapname, k))
	if err != nil {
		logger.Error("GET error: ", err.Error())
		return ""
	}
	return r
}

// func HSet(mapname, k, v string) {
// 	c := RedisClientGet().Get()
// 	defer c.Close()
// 	_, err := c.Do("HSet", mapname, k, v)
// 	if err != nil {
// 		fmt.Println("Hset error: ", err.Error())
// 		return
// 	}
// }

func HSet(mapname, k, v interface{}) {
	c := RedisClientGet().Get()
	defer c.Close()
	_, err := c.Do("HSet", mapname, k, v)
	if err != nil {
		logger.Error("Hset error: ", err.Error())
		return
	}
}

func HGetBytes(mapname, k []byte) []byte {
	c := RedisClientGet().Get()
	defer c.Close()
	r, err := redis.Bytes(c.Do("HGet", mapname, k))
	if err != nil {
		logger.Error("HGet error: ", err.Error())
		return nil
	}
	return r
}

func GetBytes(k []byte) []byte {
	c := RedisClientGet().Get()
	defer c.Close()
	r, err := redis.Bytes(c.Do("Get", k))
	if err != nil {
		logger.Error("Get error: ", err.Error())
		return nil
	}
	return r
}

func HSetBytes(mapname, k, v []byte) {
	c := RedisClientGet().Get()
	defer c.Close()
	_, err := c.Do("HSet", mapname, k, v)
	if err != nil {
		logger.Error("HSet error: ", err.Error())
		return
	}
}

func GetStringValue(k string) string {
	c := RedisClientGet().Get()
	defer c.Close()
	username, err := redis.String(c.Do("GET", k))
	if err != nil {
		logger.Error("Get Error: ", err.Error())
		return ""
	}
	return username
}

// func SetKeyExpire(k string, ex int) {
// 	c := RedisClientGet().Get()
// 	defer c.Close()
// 	_, err := c.Do("EXPIRE", k, ex)
// 	if err != nil {
// 		fmt.Println("set error", err.Error())
// 	}
// }

func SetkeyExPrire(k interface{}, ex int) {
	c := RedisClientGet().Get()
	defer c.Close()
	_, err := c.Do("EXPIRE", k, ex)
	if err != nil {
		logger.Error("set error", err.Error())
	}
}

func CheckKey(k string) bool {
	c := RedisClientGet().Get()
	defer c.Close()
	exist, err := redis.Bool(c.Do("EXISTS", k))
	if err != nil {
		logger.Error("check key failed", err)
		return false
	} else {
		return exist
	}
}

func DelKey(k string) error {
	c := RedisClientGet().Get()
	defer c.Close()
	_, err := c.Do("DEL", k)
	if err != nil {
		logger.Error("delete redis failed", err)
		return err
	}
	return nil
}

func Delete(k []byte) error {
	c := RedisClientGet().Get()
	defer c.Close()
	_, err := c.Do("DEL", k)
	if err != nil {
		logger.Error("redis delete failed", err)
		return err
	}
	return nil
}

func SetJson(k string, data interface{}) error {
	c := RedisClientGet().Get()
	defer c.Close()
	value, _ := json.Marshal(data)
	n, _ := c.Do("SETNX", k, value)
	if n != int64(1) {
		return errors.New("set failed")
	}
	return nil
}

func getJsonByte(k string) ([]byte, error) {
	c := RedisClientGet().Get()
	defer c.Close()
	jsonGet, err := redis.Bytes(c.Do("GET", k))
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	return jsonGet, nil
}

func LPush(k, v string) {
	c := RedisClientGet().Get()
	defer c.Close()
	_, err := c.Do("lpush", redis.Args{}.Add(v).AddFlat(k)...)
	if err != nil {
		logger.Error("redis LPUsh failed", err)
	}
}

func Sadd(k, v []byte) {
	c := RedisClientGet().Get()
	defer c.Close()
	_, err := c.Do("SADD", k, v)
	if err != nil {
		logger.Error("redis sadd failed,", err)
	}
}

func Geoadd(k []byte, longitude, latitude float64, member []byte) {
	c := RedisClientGet().Get()
	defer c.Close()
	_, err := c.Do("geoadd", k, longitude, latitude, member)
	if err != nil {
		logger.Error("redis geoadd failed", err)
	}
}

func Georadius(key []byte, longitude, latitude, radius float64, unit string) []vo.GeoRadius {
	c := RedisClientGet().Get()
	defer c.Close()
	result, err := c.Do("georadius", key, longitude, latitude, radius, unit, "WITHDIST", "WITHCOORD", "ASC")
	if err != nil {
		logger.Error("redis georadius failed", err)
	}
	switch v := result.(type) {
	case []interface{}:
		return parses(v)
	default:
		return nil
	}
}

func parses(data []interface{}) []vo.GeoRadius {
	var list []vo.GeoRadius
	// a := goredis.Georadius("cityGeo", 116.405285, 39.904989, 100000.0001, "km")
	for _, v := range data {
		switch n := v.(type) {
		case []interface{}:
			g := vo.GeoRadius{
				Member:        []byte(n[0].(string)),
				Distance:      parse.StringToFloat64(n[1].(string)),
				GeoCoordinate: parseGeoCoordinate(n[2]),
			}
			list = append(list, g)
		default:

		}
	}
	return list
}

func parseGeoCoordinate(data interface{}) vo.GeoCoordinate {
	switch d := data.(type) {
	case []interface{}:
		if len(d) > 1 {
			return vo.GeoCoordinate{
				Longitude: parse.StringToFloat64(d[0].(string)),
				Latitude:  parse.StringToFloat64(d[1].(string)),
			}
		}
	}
	return vo.GeoCoordinate{}
}

func Geodist(k, member1, member2 []byte, unit string) float64 {
	c := RedisClientGet().Get()
	defer c.Close()
	result, err := redis.Float64(c.Do("geodist", k, member1, member2, unit))
	if err != nil {
		logger.Error(" redis Geodist failed", err)
		return 0
	}
	return result
}

func Sinter(keys ...interface{}) []string {
	c := RedisClientGet().Get()
	defer c.Close()
	result, err := redis.Strings(c.Do("sinter", keys...))
	if err != nil {
		logger.Error("redis sinter failed", err)
		return nil
	}
	return result
}
