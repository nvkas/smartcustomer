package core

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/lukebryanshehao/smartcustomer/utils"
	"time"
)

//场景：动态多Rediser()，每个Rediser为自定义数量并发
//要求:Redis要求5.0以上
type RedisCustomers struct {
	maxDataGetCount int64           //一次获取(处理)多少数据
	stream          string          //队列名
	streamBlockTime time.Duration   //队列获取间隔
	redisGroup      []*redis.Client //redis
	//Func            func(...interface{}) //自定义消费方法
}

//新建Redis并发
//maxRunCount 并发量
//f 自定义消费方法
func NewRedisCustomers(streamName string, maxDataGetCount int64, streamBlockTime time.Duration, RedisConfigs []utils.RedisConfig) (RedisCustomers, error) {

	redisCustomers := RedisCustomers{}

	if maxDataGetCount <= 0 {
		maxDataGetCount = 1
	}
	redisCustomers.stream = streamName
	redisCustomers.streamBlockTime = streamBlockTime
	redisGroups, err := utils.InitRedis(RedisConfigs)
	if err != nil {
		return redisCustomers, err
	}
	redisCustomers.redisGroup = redisGroups
	redisCustomers.maxDataGetCount = maxDataGetCount
	//redisCustomers.Func = f

	return redisCustomers, nil
}

func (r *RedisCustomers) SetData(data map[string]interface{}) error {
	err := r.redisGroup[0].XAdd(context.TODO(), &redis.XAddArgs{
		Stream: r.stream,
		ID:     "*",
		Values: data,
	}).Err()
	return err
}

func (r *RedisCustomers) GetData() []map[string]interface{} {
	res, _ := r.redisGroup[0].XRead(context.TODO(), &redis.XReadArgs{
		Streams: []string{r.stream, "0"},
		Count:   r.maxDataGetCount,
		Block:   r.streamBlockTime,
	}).Result()

	result := make([]map[string]interface{}, 0)

	if len(res) > 0 {
		Messages := res[0].Messages
		for _, XMessage := range Messages {
			ID := XMessage.ID
			maps := XMessage.Values
			result = append(result, maps)
			r.redisGroup[0].XDel(context.TODO(), r.stream, ID)
		}
	}
	return result
}
