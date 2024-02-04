package tool

import "github.com/haiyanghan/tiangong/common/lock"

type LinkedList struct {
	lLock lock.Rwlock
	head  *Element
}

type Element struct {
	val  interface{}
	next *Element
}

func (l *LinkedList) Pop() interface{} {
	l.lLock.Lock()
	defer l.lLock.Unlock()

	if !l.Empty() {
		res := l.head.next
		l.head.next = res.next
		return res.val
	}
	return nil
}

func (l *LinkedList) Put(res interface{}) {
	l.lLock.Lock()
	defer l.lLock.Unlock()

	if l.head == nil {
		l.head = &Element{nil, nil}
	}

	if res != nil {
		t := l.head
		for t.next != nil {
			t = t.next
		}
		t.next = &Element{res, nil}
	}
}

func (l *LinkedList) Empty() bool {
	return l.head == nil || l.head.next == nil
}

func (l *LinkedList) Len() (c int) {
	l.lLock.RLock()
	defer l.lLock.RUnlock()

	if l.Empty() {
		return
	}
	h := l.head
	for h.next != nil {
		c += 1
		h = h.next
	}
	return
}
