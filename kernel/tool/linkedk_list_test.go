package tool_test

import (
	"testing"

	"github.com/haiyanghan/tiangong/kernel/tool"
)

func TestLinkedList(t *testing.T) {
	link := tool.LinkedList{}
	if !link.Empty() {
		t.Error("test fial")
		return
	}
	for i := 0; i < 10; i++ {
		a := i
		link.Put(&a)
	}
	if link.Empty() || link.Len() != 10 {
		t.Error("test fial")
		return
	}

	for i := 1; i <= 10; i++ {
		_ = link.Pop()
		if link.Len() != 10-i {
			t.Error("test fial")
			return
		}
	}
}
