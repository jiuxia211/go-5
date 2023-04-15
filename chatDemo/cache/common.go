package cache

import (
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
	"gopkg.in/ini.v1"
)

var (
	RedisClient *redis.Client
	RedisDb     string
	RedisAddr   string
	RedisPw     string
	RedisDbName string
)

func init() {
	file, err := ini.Load("./conf/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误", err)
	}
	RedisDb = file.Section("redis").Key("RedisDb").String()
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPw = file.Section("redis").Key("RedisPw ").String()
	RedisDbName = file.Section("redis").Key("RedisDbName").String()
	Redis()
}
func Redis() {
	db, _ := strconv.ParseUint(RedisDbName, 10, 64) //转化为unit64
	client := redis.NewClient(&redis.Options{
		Addr: RedisAddr,
		DB:   int(db),
	})
	//心跳测试
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	RedisClient = client
	fmt.Println("Redis 连接成功")
}
