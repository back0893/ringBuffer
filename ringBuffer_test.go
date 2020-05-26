package ringbuffer

import (
	"fmt"
	"io"
	"testing"
)

func TestRingBuffer_interface(t *testing.T) {
	rb := NewRingBuffer(1)
	var _ io.Writer = rb
	var _ io.Reader = rb
}

func TestRingBuffer_Write(t *testing.T) {
	rb := NewRingBuffer(2)
	n, err := rb.Write([]byte{1, 2})
	if err != nil {
		t.Fatal(err)
	}
	if n != 2 {
		t.Fatal(err)
	}
	if !rb.IsFull() {
		t.Fatal("计算慢错误")
	}
	_, err = rb.Write([]byte{1, 2})
	if err == nil {
		t.Fatal("幔子继续插入")
	}

}

func TestRingBuffer_Read(t *testing.T) {
	rb := NewRingBuffer(2)
	_, _ = rb.Write([]byte{1, 2})
	p := make([]byte, 1)
	_, _ = rb.Read(p)
	fmt.Println(p)
	_, _ = rb.Write([]byte{3})
	p = make([]byte, 2)
	_, _ = rb.Read(p)
	fmt.Println(p)
	_, err := rb.Read(p)
	if err == nil {
		t.Fatal("继续读取")
	}
}
