package cache

import (
	"fmt"
	"time"

	"github.com/beego/beego/v2/core/config"
	"github.com/gomodule/redigo/redis"
)

var (
	RedisCacheClient *redis.Pool
)

func init() {
	// 从配置文件获取 redis 的 ip 以及 db
	redisHost, _ := config.String("cache_redis_host")
	redisPort, _ := config.Int("cache_redis_port")
	redisDb, _ := config.Int("cache_redis_db")
	address := fmt.Sprintf("%s:%d", redisHost, redisPort)

	// 建立连接池
	RedisCacheClient = &redis.Pool{
		// 从配置文件获取 maxidle 以及 maxactive，取不到则用后面的默认值
		MaxIdle:     config.DefaultInt("cache_redis_maxidle", 16),
		MaxActive:   config.DefaultInt("cache_redis_maxactive", 512),
		IdleTimeout: 180 * time.Second,

		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", address)
			if err != nil {
				return nil, err
			}
			// 选择db
			c.Do("SELECT", redisDb)
			return c, nil
		},
	}
}
