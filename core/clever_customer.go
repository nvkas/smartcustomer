/*
@Time : 2019-12-05 14:22
@Author : Lukebryan
*/
package core

import (
	"errors"
	"sync"
)

//场景：动态多Clever()，每个Clever为自定义数量并发
type CleverCustomer struct {
	MaxUseCount int			//最大容量
	Mutex sync.Mutex		//锁
	TimeOut int	//超时不用销毁，0表示不销毁，单位：秒
	Func func(interface{})	//自定义消费方法
	CleversMap map[string]SmartCustomer	//Clever池
}

//新建并发协程
//maxRunCount 并发量
//f 自定义消费方法
func NewCleverCustomer(maxUseCount,timeOut int,f func(interface{})) CleverCustomer {
	cleverCustomer := CleverCustomer{}
	if maxUseCount <= 0 {
		maxUseCount = 10
	}
	cleverCustomer.MaxUseCount = maxUseCount
	cleverCustomer.TimeOut = timeOut
	cleverCustomer.Func = f

	maps := make(map[string]SmartCustomer)

	cleverCustomer.CleversMap = maps

	return cleverCustomer
}


//新增Clever,默认单Clever
func (c *CleverCustomer)NewClever(key string,smartRunCount int) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()

	if len(c.CleversMap) >= c.MaxUseCount {
		return errors.New("can not add more clever")
	}
	
	//默认单Clever
	if smartRunCount <= 0 {
		smartRunCount = 1
	}
	
	smartCustomer := NewSmartCustomer(smartRunCount,c.Func)
	c.CleversMap[key] = smartCustomer
	return nil
}

////获取Clever
//func (c *CleverCustomer)GetCleverSmart(key string) SmartCustomer {
//	c.Mutex.Lock()
//	defer c.Mutex.Unlock()
//	
//	return c.CleversMap[key]
//}

//销毁增Clever
func (c *CleverCustomer)Destroy(key string) {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	if _,ok := c.CleversMap[key];ok {
		smart := c.CleversMap[key]
		smart.Stop()
		delete(c.CleversMap, key)
	}
}

//获取Clever数量
func (c *CleverCustomer)GetCleverSize() int {
	//c.Mutex.Lock()
	//defer c.Mutex.Unlock()
	return len(c.CleversMap)
}



//添加数据至队列(数据入口)
func (c *CleverCustomer)AddSmartData(key string,data interface{}) {
	if _,ok := c.CleversMap[key];!ok {
		return
	}
	smart := c.CleversMap[key]
	smart.AddDataQueue(data)
	return
}

//获取数据堆积量
func (c *CleverCustomer)GetSmartDataSize(key string) int {
	if _,ok := c.CleversMap[key];!ok {
		return 0
	}
	smart := c.CleversMap[key]
	return smart.GetDataQueueSize()
}



