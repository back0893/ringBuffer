package pool

import "testing"

func TestRingBufferPool(t *testing.T) {
	rb1 := Get()
	rb1.Write([]byte{'1', '2', '3'})
	Put(rb1)
	rb2 := Get()
	if rb2.Length() != 3 {
		t.Fatal("错误")
	}
	rb3 := Get()
	if rb3.Length() != 0 {
		t.Fatal("错误")
	}
	Put(rb2)
	Put(rb3)
}
