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

	//smartCustomer
	smartCustomer := core.NewSmartCustomer(3,PrintData)
	go func() {
		i := 0
		for {
			i ++
			smartCustomer.AddDataQueue("这是"+cast.ToString(i))
			time.Sleep(time.Second*1)
		}
	}()


	//cleverCustomer
	cleverCustomer2 := core.NewCleverCustomer(10, 0, PrintData)

	var err error
	err = cleverCustomer2.NewClever("no_1", 1,nil)
	if err != nil {
		log.Println(err)
	}
	err = cleverCustomer2.NewClever("no_2", 2,PrintData2)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("size: ", cleverCustomer2.GetCleverSize())

	go func() {
		i := 0
		for {
			i ++
			cleverCustomer2.AddSmartData("no_1", "这是no_1---"+cast.ToString(i))
			time.Sleep(time.Second * 1)
			if i == 5 {
				time.Sleep(time.Second * 12)
			}
		}
	}()
	go func() {
		i := 0
		for {
			i ++
			cleverCustomer2.AddSmartData("no_2", "这是no_2---"+cast.ToString(i))
			time.Sleep(time.Second * 1)
			if i == 5 {
				time.Sleep(time.Second * 12)
			}
		}
	}()

	//smartCustomers
	smartCustomers := core.NewSmartCustomers(3,5,PrintData3)

	go func() {
		i := 0
		for {
			i ++
			smartCustomers.AddDataQueues([]interface{}{i})
			//time.Sleep(time.Second*1)
			if i == 300 {
				time.Sleep(time.Second*5)
				smartCustomers.Stop()
			}
		}
	}()

	//cleverCustomers
	cleverCustomers := core.NewCleverCustomers(10, 5,0, PrintData3)

	err = cleverCustomers.NewClevers("no_1", 5,1,nil)
	if err != nil {
		log.Println(err)
	}
	err = cleverCustomers.NewClevers("no_2", 5,2,PrintData3)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("size: ", cleverCustomers.GetCleverSize())

	go func() {
		i := 0
		for {
			i ++
			cleverCustomers.AddSmartDatas("no_1", []interface{}{"这是no_1---"+cast.ToString(i)})
			time.Sleep(time.Second * 1)
			if i == 5 {
				time.Sleep(time.Second * 12)
			}
		}
	}()
	go func() {
		i := 0
		for {
			i ++
			cleverCustomers.AddSmartDatas("no_1", []interface{}{"这是no_1---"+cast.ToString(i)})
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
func PrintData2(data interface{}) {
	fmt.Println("data--",data)
	//time.Sleep(time.Second*3)
}
func PrintData3(datas []interface{}) {
	fmt.Println("datas--",datas)
	//time.Sleep(time.Second*3)
}
