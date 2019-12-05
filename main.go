/*
@Time : 2019-12-05 14:00
@Author : Lukebryan
*/
package main

import (
	"fmt"
	"github.com/spf13/cast"
	"smartcustomer/core"
	"time"
)

func main() {
	smartCustomer := core.NewSmartCustomer(3,PrintData)
	smartCustomer2 := core.NewSmartCustomer(2,PrintData)

	go func() {
		i := 0
		for {
			i ++
			smartCustomer.AddDataQueue(i)
			time.Sleep(time.Second*1)
		}
	}()

	go func() {
		i := 0
		for {
			i ++
			smartCustomer2.AddDataQueue("这是"+cast.ToString(i))
			time.Sleep(time.Second*1)
		}
	}()
	select {

	}

}

func PrintData(data interface{})  {
	fmt.Println(data)
	time.Sleep(time.Second*3)
}
