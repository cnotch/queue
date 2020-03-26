// Copyright (c) 2018,CAO HONGJU. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package queue_test

import (
	"math/rand"
	"testing"

	"github.com/cnotch/queue"
)

func TestQueue_Basic(t *testing.T) {
	var q queue.Queue
	q.Push(3)
	if q.Len() != 1 {
		t.Error("Failed Queue.Enqueue")
	}
}

func TestQueue_Pop(t *testing.T) {
	var q queue.Queue
	if _, ok := q.Pop(); ok {
		t.Error("Failed Queue.Top")
	}
	q.Push("test")
	q.Push(3)
	if value, _ := q.Pop(); !(value == "test" && q.Len() == 1) {
		t.Errorf("Failed Queue.Dequeue, value is %d, len is %d", value, q.Len())
	}
}

func TestQueue_Empty(t *testing.T) {
	var q queue.Queue
	if !q.Empty() {
		t.Error("Failed Queue.Empty")
	}
}

func BenchmarkQueue(b *testing.B) {
	var q queue.Queue
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		n := rand.Intn(100)
		for j := 0; j < n; j++ {
			q.Push(j)
		}
		m := rand.Intn(100)
		for j := 0; j < m; j++ {
			q.Pop()
		}
	}
}
