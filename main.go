/*
@Time : 2019-12-05 14:00
@Author : Lukebryan
*/
package main

import (
	"fmt"
	"github.com/lukebryanshehao/smartcustomer/core"
	"github.com/spf13/cast"
	"log"
	"time"
)

func main() {
	//计数器
	//counter := utils.NewCounter()
	//
	//go func() {
	//	time.Sleep(time.Second * 15)
	//	counter.ReStart()
	//	time.Sleep(time.Second * 15)
	//	counter.Stop()
	//}()
	//
	//for {
	//	fmt.Println("count: ",counter.Count)
	//	time.Sleep(time.Second * 1)
	//}

	//单方法并发
	//smartCustomer := core.NewSmartCustomer(3,PrintData)
	//smartCustomer2 := core.NewSmartCustomer(2,PrintData)
	//
	//go func() {
	//	i := 0
	//	for {
	//		i ++
	//		smartCustomer.AddDataQueue(i)
	//		time.Sleep(time.Second*1)
	//		if i == 30 {
	//			smartCustomer.Stop()
	//		}
	//	}
	//}()
	//
	//go func() {
	//	i := 0
	//	for {
	//		i ++
	//		smartCustomer2.AddDataQueue("这是"+cast.ToString(i))
	//		time.Sleep(time.Second*1)
	//	}
	//}()
	//select {
	//
	//}

	//Clever
	cleverCustomer := core.NewCleverCustomer(10, 0, PrintData)

	var err error
	err = cleverCustomer.NewClever("no_1", 1)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("size: ", cleverCustomer.GetCleverSize())

	//go func() {
	//	time.Sleep(time.Second*30)
	//	cleverCustomer.Destroy("wxid_1")
	//}()

	go func() {
		i := 0
		for {
			i ++
			cleverCustomer.AddSmartData("no_1", "这是no_1---"+cast.ToString(i))
			time.Sleep(time.Second * 1)
			if i == 5 {
				time.Sleep(time.Second * 12)
			}
		}
	}()
	select {}

}

func PrintData(data interface{}) {
	fmt.Println(data)
	//time.Sleep(time.Second*3)
}
