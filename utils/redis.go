package utils

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type RedisConfig struct {
	Host string
	Port string
	Password string
} 

func InitRedis(RedisConfigs []RedisConfig) ([]*redis.Client,error) {
	var redisGroup []*redis.Client
	for _, RedisConfig := range RedisConfigs {
		conn,err := openConn(RedisConfig.Host,RedisConfig.Port,RedisConfig.Password)
		if err != nil {
			return redisGroup,err
		}
		redisGroup = append(redisGroup, conn)
	}
	return redisGroup,nil
}

func openConn(Host,Port,Password string) (*redis.Client,error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", Host, Port),
		Password: Password, // 密码
	})
	if _, err := client.Ping(context.Background()).Result(); err != nil {
		return nil,err
	}
	return client,nil
}