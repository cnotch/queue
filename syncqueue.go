// Copyright (c) 2019,CAOHONGJU All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package queue

import (
	"sync"
)

// SyncQueue 同步队列
type SyncQueue struct {
	cond  *sync.Cond
	queue Queue
}

// NewSyncQueue 创建同步队列
func NewSyncQueue() *SyncQueue {
	return &SyncQueue{
		cond: sync.NewCond(&sync.Mutex{}),
	}
}

// Queue 返回内部队列
func (q *SyncQueue) Queue() *Queue {
	return &q.queue
}

// Push 入列并发送信号
func (q *SyncQueue) Push(e Elem) {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	q.queue.Push(e)
	q.cond.Signal()
}

// Pop 出列，如果没有等待信号做一次重试
func (q *SyncQueue) Pop() Elem {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()

	if q.queue.Len() <= 0 {
		q.cond.Wait()
	}

	e, _ := q.queue.Pop()
	return e
}

// Signal 发送信号，以便结束等待
func (q *SyncQueue) Signal() {
	q.cond.Signal()
}

// Broadcast 广播信号，以释放所有出列的阻塞等待
func (q *SyncQueue) Broadcast() {
	q.cond.Broadcast()
}

// Len 队列长度
func (q *SyncQueue) Len() int {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	return q.queue.Len()
}

// Reset 重置队列
func (q *SyncQueue) Reset() {
	q.cond.L.Lock()
	defer q.cond.L.Unlock()
	q.queue.Reset()
}
