/*
@Time : 2019-12-05 14:22
@Author : Lukebryan
*/
package core

import (
	"errors"
	"github.com/lukebryanshehao/smartcustomer/utils"
	"sync"
	"time"
)

//场景：动态多Clever()，每个Clever为自定义数量并发
type CleverCustomers struct {
	MaxUseCount     int                      //最大容量
	MaxDataGetCount int                      //一次获取(处理)多少数据
	Mutex           sync.Mutex               //锁
	TimeOut         int                      //超时不用销毁，0表示不销毁，单位：秒
	Func            func([]interface{})        //自定义消费方法
	SmartMaps       map[string]SmartCustomers //Clever池
	CounterMaps     map[string]*utils.Counter
}

//新建并发协程
//maxRunCount 并发量
//f 自定义消费方法
func NewCleverCustomers(maxUseCount,maxDataGetCount, timeOut int, f func([]interface{})) *CleverCustomers {
	CleverCustomers := CleverCustomers{}
	if maxUseCount <= 0 {
		maxUseCount = 10
	}
	if maxDataGetCount <= 0 {
		maxDataGetCount = 1
	}

	CleverCustomers.MaxUseCount = maxUseCount
	CleverCustomers.MaxDataGetCount = maxDataGetCount
	CleverCustomers.TimeOut = timeOut
	CleverCustomers.Func = f

	maps := make(map[string]SmartCustomers)
	CleverCustomers.SmartMaps = maps

	maps2 := make(map[string]*utils.Counter)
	CleverCustomers.CounterMaps = maps2

	return &CleverCustomers
}

//新增Clever,默认单Clever
func (c *CleverCustomers) NewClevers(key string, smartRunCount,maxDataGetCount int, f func([]interface{})) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if len(c.SmartMaps) >= c.MaxUseCount {
		return errors.New("can not add more clever")
	}

	//默认单Clever
	if smartRunCount <= 0 {
		smartRunCount = 1
	}
	dataGetCount := c.MaxDataGetCount
	if maxDataGetCount > 0 {
		dataGetCount = maxDataGetCount
	}
	if dataGetCount <= 0 {
		dataGetCount = 1
	}
	fun := c.Func
	if f != nil {
		fun = f
	}
	smartCustomer := NewSmartCustomers(smartRunCount,dataGetCount, fun)
	c.SmartMaps[key] = smartCustomer

	counter := utils.NewCounter()
	c.CounterMaps[key] = counter

	if c.TimeOut > 0 {
		go c.CheckTimeOut(key)
	}
	return nil
}

//销毁增Clever
func (c *CleverCustomers) Destroy(key string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	if _, ok := c.SmartMaps[key]; ok {
		smart := c.SmartMaps[key]
		smart.Stop()
		delete(c.SmartMaps, key)
	}
	if _, ok := c.CounterMaps[key]; ok {
		counter := c.CounterMaps[key]
		counter.Stop()
		delete(c.CounterMaps, key)
	}
}

//获取Clever数量
func (c *CleverCustomers) GetCleverSize() int {
	//c.Mutex.Lock()
	//defer c.Mutex.Unlock()
	return len(c.SmartMaps)
}

//添加数据至队列(数据入口)
func (c *CleverCustomers) AddSmartDatas(key string, datas []interface{}) {
	if _, ok := c.SmartMaps[key]; !ok {
		return
	}
	smart := c.SmartMaps[key]
	smart.AddDataQueues(datas)

	c.CounterMaps[key].ReStart()

	return
}

//获取数据堆积量
func (c *CleverCustomers) GetSmartDataSize(key string) int {
	if _, ok := c.SmartMaps[key]; !ok {
		return 0
	}
	smart := c.SmartMaps[key]
	return smart.GetDataQueueSize()
}

func (c *CleverCustomers) CheckTimeOut(key string) {

	for {
		if _, ok := c.SmartMaps[key]; !ok {
			return
		}

		count := c.CounterMaps[key].Count

		if count >= c.TimeOut {
			c.Destroy(key)
			return
		}

		time.Sleep(time.Second * 1)
	}
}
