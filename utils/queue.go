package utils

import "sync"

type Item interface {}

// 队列结构体
type Queue struct {
	Items []Item
}

// 新建
func (q *Queue)New() *Queue  {
	q.Items = []Item{}
	return q
}
var mutex sync.Mutex
// 入队
func (q *Queue) Enqueue(data Item)  {
	mutex.Lock()
	q.Items = append(q.Items, data)
	mutex.Unlock()
}

// 出队
func (q *Queue) Dequeue() *Item {
	mutex.Lock()
	defer mutex.Unlock()
	// 由于是先进先出，0为队首
	if q.IsEmpty() {
		return nil
	}
	item := q.Items[0]
	q.Items = q.Items[1: len(q.Items)]

	return &item
}

// 队列是否为空
func (q *Queue) IsEmpty() bool  {
	return len(q.Items) == 0
}

// 队列置为空
func (q *Queue) ToEmpty() {
	mutex.Lock()
	defer mutex.Unlock()
	q.Items = []Item{}
}

// 队列长度
func (q *Queue) Size() int  {
	return len(q.Items)
}

func NewQueue() *Queue  {
	q := Queue{}
	queue := q.New()
	return queue
}
