/*
@Time : 2020-04-10 10:39
@Author : Lukebryan
*/
package utils

import (
	"time"
)

type Counter struct {
	StopFlag bool
	Count int
}

func NewCounter() *Counter {
	counter := Counter{}

	go counter.Start()

	return &counter
}

func (c *Counter)Start() {

	for {

		if c.StopFlag == true {
			return
		}

		c.Count = c.Count + 1
		time.Sleep(time.Second * 1)
	}

}

func (c *Counter)ReStart() {
	c.Count = 0
}

func (c *Counter)Stop() {
	c.StopFlag = true
}
