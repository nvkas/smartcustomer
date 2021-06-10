/*
@Time : 2019-12-05 14:22
@Author : Lukebryan
*/
package core

import (
	"github.com/lukebryanshehao/smartcustomer/utils"
	"sync"
	"time"
)

type SmartCustomers struct {
	MaxRunCount     int                 //最大并发
	MaxDataGetCount int                 //一次获取(处理)多少数据
	DataQueue       *utils.Queue        //业务数据
	Mutex           sync.Mutex          //锁
	Func            func([]interface{}) //自定义消费方法
	runningFlag bool	//运行标志
}

//新建并发协程
//maxRunCount 并发量
//f 自定义消费方法
func NewSmartCustomers(maxRunCount, maxDataGetCount int, f func([]interface{})) SmartCustomers {
	smartCustomers := SmartCustomers{}
	if maxRunCount <= 0 {
		maxRunCount = 50
	}
	if maxDataGetCount <= 0 {
		maxDataGetCount = 1
	}
	smartCustomers.MaxRunCount = maxRunCount
	smartCustomers.MaxDataGetCount = maxDataGetCount
	smartCustomers.DataQueue = utils.NewQueue()
	smartCustomers.Func = f
	smartCustomers.runningFlag = true

	go smartCustomers.queueCustomer()

	return smartCustomers
}

func (s *SmartCustomers) Stop() {
	s.DataQueue.ToEmpty()
	s.runningFlag = false
}

//添加数据至队列(数据入口)
func (s *SmartCustomers) AddDataQueues(datas []interface{}) {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	if !s.runningFlag {
		return
	}
	for _, data := range datas {
		s.DataQueue.Enqueue(data)
	}
	return
}

//从队列拿数据
func (s *SmartCustomers) getDataQueues() []interface{} {
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	getCount := s.MaxDataGetCount
	dataSize := s.GetDataQueueSize()
	if getCount > dataSize {
		getCount = dataSize
	}
	var datas []interface{}
	for i := 0; i < getCount; i++ {
		data := s.DataQueue.Dequeue()
		if data == nil {
			continue
		}
		datas = append(datas, *data)
	}

	return datas
}

//获取数据堆积量
func (s *SmartCustomers) GetDataQueueSize() int {
	v := s.DataQueue.Size()
	return v
}

//消费
func (s *SmartCustomers) queueCustomer() {
	ch := make(chan int, s.MaxRunCount)
	wg := sync.WaitGroup{}
	for {
		datas := s.getDataQueues()

		if !s.runningFlag {
			return
		}

		ch <- 1
		wg.Add(1)
		go s.dataCustomers(&wg, ch, datas)
	}
	wg.Wait()
}

func (s *SmartCustomers) dataCustomers(wg *sync.WaitGroup, ch chan int, datas []interface{}) {
	defer func() {
		wg.Done()
		<-ch
	}()

	if len(datas) == 0 {
		time.Sleep(time.Second * 1)
		return
	}

	//TODO 你要消费数据的业务
	s.Func(datas)
}
