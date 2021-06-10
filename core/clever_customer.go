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
type CleverCustomer struct {
	MaxUseCount int			//最大容量
	Mutex sync.Mutex		//锁
	TimeOut int	//超时不用销毁，0表示不销毁，单位：秒
	Func func(interface{})	//自定义消费方法
	SmartMaps map[string]SmartCustomer	//Clever池
	CounterMaps map[string]*utils.Counter
}

//新建并发协程
//maxRunCount 并发量
//f 自定义消费方法
func NewCleverCustomer(maxUseCount,timeOut int,f func(interface{})) CleverCustomer {
	cleverCustomer := CleverCustomer{}
	if maxUseCount <= 0 {
		maxUseCount = 10
	}
	//
	//if timeOut <= 0 {
	//	timeOut = 365 * 24 * 60 * 60
	//}

	cleverCustomer.MaxUseCount = maxUseCount
	cleverCustomer.TimeOut = timeOut
	cleverCustomer.Func = f

	maps := make(map[string]SmartCustomer)
	cleverCustomer.SmartMaps = maps

	maps2 := make(map[string]*utils.Counter)
	cleverCustomer.CounterMaps = maps2

	return cleverCustomer
}


//新增Clever,默认单Clever
func (c *CleverCustomer)NewClever(key string,smartRunCount int,f func(interface{})) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if len(c.SmartMaps) >= c.MaxUseCount {
		return errors.New("can not add more clever")
	}
	
	//默认单Clever
	if smartRunCount <= 0 {
		smartRunCount = 1
	}
	fun := c.Func
	if f != nil {
		fun = f
	}
	smartCustomer := NewSmartCustomer(smartRunCount,fun)
	c.SmartMaps[key] = smartCustomer

	counter := utils.NewCounter()
	c.CounterMaps[key] = counter

	if c.TimeOut > 0 {
		go c.CheckTimeOut(key)
	}
	return nil
}

////获取Clever
//func (c *CleverCustomer)GetCleverSmart(key string) SmartCustomer {
//	c.Mutex.Lock()
//	defer c.Mutex.Unlock()
//	
//	return c.SmartMaps[key]
//}

//销毁增Clever
func (c *CleverCustomer)Destroy(key string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	if _,ok := c.SmartMaps[key];ok {
		smart := c.SmartMaps[key]
		smart.Stop()
		delete(c.SmartMaps, key)
	}
	if _,ok := c.CounterMaps[key];ok {
		counter := c.CounterMaps[key]
		counter.Stop()
		delete(c.CounterMaps, key)
	}
}

//获取Clever数量
func (c *CleverCustomer)GetCleverSize() int {
	//c.Mutex.Lock()
	//defer c.Mutex.Unlock()
	return len(c.SmartMaps)
}



//添加数据至队列(数据入口)
func (c *CleverCustomer)AddSmartData(key string,data interface{}) {
	if _,ok := c.SmartMaps[key];!ok {
		return
	}
	smart := c.SmartMaps[key]
	smart.AddDataQueue(data)

	c.CounterMaps[key].ReStart()

	return
}

//获取数据堆积量
func (c *CleverCustomer)GetSmartDataSize(key string) int {
	if _,ok := c.SmartMaps[key];!ok {
		return 0
	}
	smart := c.SmartMaps[key]
	return smart.GetDataQueueSize()
}

func (c *CleverCustomer) CheckTimeOut(key string) {

	for {
		if _,ok := c.SmartMaps[key];!ok {
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



