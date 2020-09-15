package redis

import (
	"fmt"
	"sprout_server/settings"

	"github.com/go-redis/redis"
)

// declare a global rdb variable
var rdb *redis.Client

// Init connection
func Init(cfg *settings.RedisConfig) (err error) {
	rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
		PoolSize: cfg.PoolSize,
	})

	_, err = rdb.Ping().Result()
	return
}

// Because db is private,
// we provide the public Close for other packages
// to close the db connection
func Close() {
	_ = rdb.Close()
}
