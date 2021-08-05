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
	maxUseCount     int                      //最大容量
	maxDataGetCount int                      //一次获取(处理)多少数据
	mutex           sync.Mutex               //锁
	timeOut         int                      //超时不用销毁，0表示不销毁，单位：秒
	Func            func([]interface{})        //自定义消费方法
	smartMaps       map[string]SmartCustomers //Clever池
	counterMaps     map[string]*utils.Counter
}

//新建并发协程
//maxRunCount 并发量
//f 自定义消费方法
func NewCleverCustomers(maxUseCount,maxDataGetCount, timeOut int, f func([]interface{})) CleverCustomers {
	cleverCustomers := CleverCustomers{}
	if maxUseCount <= 0 {
		maxUseCount = 10
	}
	if maxDataGetCount <= 0 {
		maxDataGetCount = 1
	}

	cleverCustomers.maxUseCount = maxUseCount
	cleverCustomers.maxDataGetCount = maxDataGetCount
	cleverCustomers.timeOut = timeOut
	cleverCustomers.Func = f

	maps := make(map[string]SmartCustomers)
	cleverCustomers.smartMaps = maps

	maps2 := make(map[string]*utils.Counter)
	cleverCustomers.counterMaps = maps2

	return cleverCustomers
}

//新增Clever,默认单Clever
func (c *CleverCustomers) NewClevers(key string, smartRunCount,maxDataGetCount int, f func([]interface{})) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.smartMaps) >= c.maxUseCount {
		return errors.New("can not add more clever")
	}

	//默认单Clever
	if smartRunCount <= 0 {
		smartRunCount = 1
	}
	dataGetCount := c.maxDataGetCount
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
	c.smartMaps[key] = smartCustomer

	counter := utils.NewCounter()
	c.counterMaps[key] = counter

	if c.timeOut > 0 {
		go c.CheckTimeOut(key)
	}
	return nil
}

//销毁增Clever
func (c *CleverCustomers) Destroy(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	if _, ok := c.smartMaps[key]; ok {
		smart := c.smartMaps[key]
		smart.Stop()
		delete(c.smartMaps, key)
	}
	if _, ok := c.counterMaps[key]; ok {
		counter := c.counterMaps[key]
		counter.Stop()
		delete(c.counterMaps, key)
	}
}

//获取Clever数量
func (c *CleverCustomers) GetCleverSize() int {
	//c.Mutex.Lock()
	//defer c.Mutex.Unlock()
	return len(c.smartMaps)
}

//添加数据至队列(数据入口)
func (c *CleverCustomers) AddSmartDatas(key string, datas []interface{}) {
	if _, ok := c.smartMaps[key]; !ok {
		return
	}
	smart := c.smartMaps[key]
	smart.AddDataQueues(datas)

	c.counterMaps[key].ReStart()

	return
}

//获取数据堆积量
func (c *CleverCustomers) GetSmartDataSize(key string) int {
	if _, ok := c.smartMaps[key]; !ok {
		return 0
	}
	smart := c.smartMaps[key]
	return smart.GetDataQueueSize()
}

func (c *CleverCustomers) CheckTimeOut(key string) {

	for {
		if _, ok := c.smartMaps[key]; !ok {
			return
		}

		count := c.counterMaps[key].GetCount()

		if count >= c.timeOut {
			c.Destroy(key)
			return
		}

		time.Sleep(time.Second * 1)
	}
}
