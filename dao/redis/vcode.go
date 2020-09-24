package redis

import (
	"sprout_server/common/mytime"
	"time"
)

const ECode = "sprout:ecode"
const ECodeCount = "sprout:ecode:count"

func GetECode(key string) (code string, err error) {
	code, err = rdb.Get(joinKey(ECode, key)).Result()
	return
}

func GetECodeCount(key string) (count int, err error) {
	count, err = rdb.Get(joinKey(ECodeCount, key)).Int()
	return
}

func SetECode(key string, eCode string, expire time.Duration) (value string, err error) {
	value, err = rdb.Set(joinKey(ECode, key), eCode, expire).Result()
	return
}

func IncrECodeCount(key string) (count int64, err error) {
	key = joinKey(ECodeCount, key)
	count, err = rdb.Incr(key).Result()
	rdb.ExpireAt(key, mytime.TomorrowZero())
	return
}

func joinKey(prefix string, suffix string) string {
	return prefix + ":" + suffix
}
