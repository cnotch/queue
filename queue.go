// Copyright (c) 2018,CAOHONGJU. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
//
// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package queue

// Elem represents an element of a queue.
type Elem = interface{}

// Queue  represents a queue.
// The zero value for Queue is an empty queue ready to use.
type Queue struct {
	buf []Elem // contents  buf[off : len(buf)]
	off int    // read at &buf[off], write at &buf[len(buf)]
}

// NewQueue creates and initializes a new Queue using buf as its
// initial contents. The new Queue takes ownership of buf, and the
// caller should not use buf after this call. NewQueue is intended to
// prepare a Queue to read existing data. It can also be used to set
// the initial size of the internal queue for writing. To do that,
// buf should have the desired capacity but a length of zero.
//
// In most cases, new(Queue) (or just declaring a Queue variable) is
// sufficient to initialize a Queue.
func NewQueue(buf []Elem) *Queue {
	return &Queue{buf: buf}
}

// Empty determines if the queue is empty.
func (q *Queue) Empty() bool { return len(q.buf) <= q.off }

// Elems return all elements of queue.
func (q *Queue) Elems() []Elem { return q.buf[q.off:] }

// Len returns the number of elements of queue; q.Len() == len(q.Elems()).
func (q *Queue) Len() int { return len(q.buf) - q.off }

// Cap returns the capacity of the queue's underlying slice
func (q *Queue) Cap() int { return cap(q.buf) }

// Reset resets the queue to be empty, but it retains the underlying storage for use by future.
func (q *Queue) Reset() {
	resetSlice(q.buf[q.off:])

	q.buf = q.buf[:0]
	q.off = 0
}

// Grow grows the queue's capacity, if necessary, to guarantee space for
// another n elements. After Grow(n), at least n elements can be pushed to the
// queue without another allocation.
// If n is negative, Grow will panic.
func (q *Queue) Grow(n int) {
	if n < 0 {
		panic("Queue.Grow: negative count")
	}
	m := q.grow(n)
	q.buf = q.buf[:m]
}

// Get returns the element at index i in the queue.
// If the index is invalid, the call will panic.
func (q *Queue) Get(i int) (e Elem) {
	return q.buf[q.off+i]
}

// Push pushes the element e to the tail of the queue.
func (q *Queue) Push(e Elem) {
	m, ok := q.tryGrowByReslice(1)
	if !ok {
		m = q.grow(1)
	}
	q.buf[m] = e
}

// PushN pushes the contents of elems to the tail of the queue.
func (q *Queue) PushN(elems []Elem) {
	m, ok := q.tryGrowByReslice(len(elems))
	if !ok {
		m = q.grow(len(elems))
	}
	copy(q.buf[m:], elems)
}

// Pop removes and returns the front element from the queue.
// Ok is true if the queue is empty, otherwise is false.
func (q *Queue) Pop() (v Elem, ok bool) {
	if q.Empty() {
		// Queue is empty, reset to recover space.
		q.Reset()
		return
	}

	v, ok = q.buf[q.off], true
	resetSlice(q.buf[q.off : q.off+1]) // avoid memory leaks
	q.off++
	return
}

// PopN copies and removes the front len(elems) elements from
// the queue or until the queue is drained.
// The return value n is the number of elements copied.
func (q *Queue) PopN(elems []Elem) (n int) {
	if q.Empty() {
		// Queue is empty, reset to recover space.
		q.Reset()
		return
	}

	n = copy(elems, q.buf[q.off:])
	resetSlice(q.buf[q.off : q.off+n]) // avoid memory leaks

	q.off += n
	return
}

// Skip skips the front n elements.
// The return value n is the number of elements discarded.
func (q *Queue) Skip(n int) int {
	if q.Empty() {
		// Queue is empty, reset to recover space.
		q.Reset()
		return 0
	}

	len := q.Len()
	if n > len {
		n = len
	}
	resetSlice(q.buf[q.off : q.off+n]) // avoid memory leaks

	q.off += n
	return n
}

func resetSlice(sl []Elem) {
	var nilValue Elem
	for i := 0; i < len(sl); i++ {
		sl[i] = nilValue // Notify the GC early to avoid memory leaks
	}
}

func (q *Queue) tryGrowByReslice(n int) (int, bool) {
	if l := len(q.buf); n <= cap(q.buf)-l {
		q.buf = q.buf[:l+n]
		return l, true
	}
	return 0, false
}

func (q *Queue) grow(n int) int {
	m := q.Len()
	// If buffer is empty, reset to recover space.
	if m == 0 && q.off != 0 {
		q.Reset()
	}

	// Try to grow by means of a reslice.
	if i, ok := q.tryGrowByReslice(n); ok {
		return i
	}

	c := cap(q.buf)
	if n <= c/2-m {
		// We can slide things down instead of allocating a new
		// slice. We only need m+n <= c to slide, but
		// we instead let capacity get twice as large so we
		// don't spend all our time copying.
		copy(q.buf, q.buf[q.off:])
		resetSlice(q.buf[m:]) // avoid memory leaks
	} else {
		// Not enough space anywhere, we need to allocate.
		buf := make([]Elem, 2*c + n)
		copy(buf, q.buf[q.off:])
		resetSlice(q.buf[q.off:]) // avoid memory leaks
		q.buf = buf
	}

	// set new off and buf
	q.off = 0
	q.buf = q.buf[:m+n]
	return m
}
