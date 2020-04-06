/*
@Time : 2019-12-05 14:00
@Author : Lukebryan
*/
package main

import (
	"fmt"
	"github.com/spf13/cast"
	"log"
	"smartcustomer/core"
	"time"
)

func main() {
	//smartCustomer := core.NewSmartCustomer(3,PrintData)
	//smartCustomer2 := core.NewSmartCustomer(2,PrintData)
	//
	//go func() {
	//	i := 0
	//	for {
	//		i ++
	//		smartCustomer.AddDataQueue(i)
	//		time.Sleep(time.Second*1)
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

	cleverCustomer := core.NewCleverCustomer(10,0,PrintData)

	var err error
	err = cleverCustomer.NewClever("wxid_1",1)
	err = cleverCustomer.NewClever("wxid_2",1)
	if err != nil {
		log.Println(err)
	}

	fmt.Println("size: ",cleverCustomer.GetCleverSize())

	go func() {
		time.Sleep(time.Second*30)
		cleverCustomer.DestroyClever("wxid_1")
	}()

	go func() {
		i := 0
		for {
			i ++
			cleverCustomer.AddSmartData("wxid_1","这是wxid_1---"+cast.ToString(i))
			time.Sleep(time.Second*1)
			fmt.Println("smart data size: ",cleverCustomer.GetSmartDataSize("wxid_1"))
		}
	}()
	select {

	}

}

func PrintData(data interface{})  {
	fmt.Println(data)
	time.Sleep(time.Second*3)
}
