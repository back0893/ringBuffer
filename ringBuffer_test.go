package ringbuffer

import (
	"fmt"
	"io"
	"strings"
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
func TestRingBuffer_Wr(t *testing.T) {
	rb := NewRingBuffer(2)
	if err := rb.WriteByte('1'); err != nil {
		t.Fatal(err)
	}
	if _, err := rb.ReadByte(); err != nil {
		t.Fatal(err)
	}
	if err := rb.WriteByte('2'); err != nil {
		t.Fatal(err)
	}
	rb.WriteByte('2')
	if rb.Length() != 2 {
		t.Fatal("长度错误")
	}
	if rb.Free() != 0 {
		t.Fatal("长度错误")
	}
	if !rb.IsFull() {
		t.Fatal("长度错误")
	}
}
func TestRindBufferVirtualRead(t *testing.T) {
	bf := NewRingBuffer(3)
	bf.WriteByte('1')
	bf.WriteByte('2')
	bf.WriteByte('3')
	target := make([]byte, 3)
	if n, err := bf.VirtualRead(target); err != nil {
		t.Fatal(err)
	} else if n != 3 {
		t.Fatal("读取错误")
	}
	if bf.r != 0 {
		t.Fatal("虚拟错误读取")
	}
	bf.Read(target)
	if bf.r != 0 {
		t.Fatal("错误读取")
	}
	if _, err := bf.Read(target); err == nil {
		t.Fatal("错误读取")
	}
	t.Log(bf.Length(), bf.Free())
}
func TestRindBufferIndex(t *testing.T) {
	bf := NewRingBuffer(6)
	bf.Write([]byte("1p23mm"))
	bf.ReadByte()
	bf.WriteByte('p')
	i := bf.Index([]byte{'p'}, 0)
	t.Log(i)
	j := bf.Index([]byte{'p'}, i+1)
	t.Log(i, j)
	bf.Peek(i)
	bh := make([]byte, j-i+1)
	n, err := bf.Read(bh)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(bh[:n]))
}
func TestRingBufferAutoSize(t *testing.T) {
	rb := NewRingBuffer(2)
	rb.SetMaxSize(10)
	rb.SetAutoSize(2)
	rb.Write([]byte{'1', '2'})
	rb.Write([]byte{'1', '2'})
	if rb.Capacity() != 4 {
		t.Fatal("扩容失败", rb.Capacity())
	}
	if _, err := rb.Write([]byte(strings.Repeat("a", 6))); err != nil {
		t.Fatal("扩容成功")
	}

	if _, err := rb.Write([]byte{'1', '2'}); err == nil {
		t.Fatal("扩容成功")
	}

}
