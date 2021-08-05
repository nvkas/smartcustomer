/*
@Time : 2019-12-05 14:22
@Author : Lukebryan
*/
package core

import (
	"github.com/lukebryanshehao/smartcustomer/utils"
	"github.com/spf13/cast"
	"sync"
	"time"
)

type SmartCustomer struct {
	maxRunCount int               //最大并发
	dataQueue   *utils.Queue      //业务数据
	mutex       sync.Mutex        //锁
	Func        func(interface{}) //自定义消费方法
}

//新建并发协程
//maxRunCount 并发量
//f 自定义消费方法
func NewSmartCustomer(maxRunCount int, f func(interface{})) SmartCustomer {
	smartCustomer := SmartCustomer{}
	if maxRunCount <= 0 {
		maxRunCount = 50
	}
	smartCustomer.maxRunCount = maxRunCount
	smartCustomer.dataQueue = utils.NewQueue()
	smartCustomer.Func = f

	go smartCustomer.queueCustomer()

	return smartCustomer
}

func (s *SmartCustomer) Stop() {
	s.dataQueue.ToEmpty()
	s.AddDataQueue("stopSmartCustomer")
}

//添加数据至队列(数据入口)
func (s *SmartCustomer) AddDataQueue(data interface{}) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.dataQueue.Enqueue(data)
	return
}

//从队列拿数据
func (s *SmartCustomer) getDataQueue() interface{} {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	data := s.dataQueue.Dequeue()
	if data == nil {
		return nil
	}
	return *data
}

//获取数据堆积量
func (s *SmartCustomer) GetDataQueueSize() int {
	v := s.dataQueue.Size()
	return v
}

//消费
func (s *SmartCustomer) queueCustomer() {
	ch := make(chan int, s.maxRunCount)
	wg := sync.WaitGroup{}
	for {
		data := s.getDataQueue()

		if cast.ToString(data) == "stopSmartCustomer" {
			return
		}

		ch <- 1
		wg.Add(1)
		go s.dataCustomer(&wg, ch, data)
	}
	wg.Wait()
}

func (s *SmartCustomer) dataCustomer(wg *sync.WaitGroup, ch chan int, data interface{}) {
	defer func() {
		wg.Done()
		<-ch
	}()

	if data == nil {
		time.Sleep(time.Second * 1)
		return
	}

	//TODO 你要消费数据的业务
	s.Func(data)
}
