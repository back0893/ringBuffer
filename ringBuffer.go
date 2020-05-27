package ringbuffer

import (
	"errors"
)

var (
	errFull  = errors.New("不能新增值")
	errEmpty = errors.New("空")
)

type RingBuffer struct {
	data    []byte
	r       int //下一个读取的位置
	w       int //下一个写入的位置
	size    int
	isEmpty bool
}

func NewRingBuffer(length int) *RingBuffer {
	return &RingBuffer{
		data:    make([]byte, length),
		r:       0,
		w:       0,
		size:    length,
		isEmpty: true,
	}
}
func (r *RingBuffer) Read(target []byte) (int, error) {
	n := len(target)
	if n == 0 {
		return 0, nil
	}
	valiLen := r.Length()
	if valiLen == 0 {
		return 0, errEmpty
	}
	if n > valiLen {
		n = valiLen
	}
	if r.r < r.w {
		copy(target, r.data[r.r:r.w])
		r.r += n
	} else {
		c1 := r.size - r.r
		if c1 >= n {
			copy(target, r.data[r.r:])
			r.r += n
		} else {
			copy(target[:c1], r.data[r.r:])
			c2 := n - r.r
			copy(target[c1:], r.data[:c2])
			r.r = c2
		}
	}
	if r.r == r.size {
		r.r = 0
	}
	r.isEmpty = true
	return n, nil
}
func (r *RingBuffer) ReadByte() (byte, error) {
	if r.Length() == 0 {
		return 0, errEmpty
	}
	b := r.data[r.r]
	r.r += 1
	if r.r == r.size {
		r.r = 0
	}
	r.isEmpty = true
	return b, nil
}
func (r *RingBuffer) Write(data []byte) (int, error) {
	if len(data) == 0 {
		return 0, nil
	}
	if r.IsFull() {
		return 0, errFull
	}
	free := r.Free()
	n := len(data)
	if n > free {
		return 0, errFull
	}
	if r.r > r.w {
		copy(r.data[r.w:], data)
		r.w += n
	} else {
		c1 := r.size - r.w
		if c1 > n {
			copy(r.data[r.w:], data)
			r.w += n
		} else {
			copy(r.data[r.w:], data[:c1])
			c2 := n - c1
			copy(r.data[0:], data[c1:])
			r.w = c2
		}

	}
	if r.w == r.size {
		r.w = 0
	}
	r.isEmpty = false
	return n, nil
}
func (r *RingBuffer) WriteByte(b byte) error {
	if r.IsFull() {
		return errFull
	}
	r.data[r.w] = b
	r.w++
	if r.w == r.size {
		r.w = 0
	}
	r.isEmpty = false
	return nil
}
func (r *RingBuffer) Capacity() int {
	return r.size
}
func (r *RingBuffer) Reset() {
	r.r = 0
	r.w = 0
	r.isEmpty = true
}

//保存的数据长度
func (r *RingBuffer) Length() int {
	return r.size - r.Free()
}

//读指针前进
func (r *RingBuffer) Retrieve(l int) {
	if l == 0 {
		return
	}
	valiLength := r.Length()
	if valiLength == 0 {
		return
	}
	if l > valiLength {
		l = valiLength
	}
	r.r = (r.r + l) % r.size
	if r.w == r.r {
		r.isEmpty = true
	}
}

//读指针获得更多的byte
func (r *RingBuffer) Peek(i int) ([]byte, error) {
	if i == 0 {
		return []byte{}, nil
	}
	data := make([]byte, i)
	if _, err := r.Write(data); err != nil {
		return nil, err
	}
	return data, nil
}

//获得当前全部可以读取的bytes,但是不会移动读指针
//只是一个拷贝,防治修改值
func (r *RingBuffer) Bytes() []byte {
	valiLenth := r.Length()
	if valiLenth == 0 {
		return []byte{}
	}
	data := make([]byte, valiLenth)
	if r.w > r.r {
		copy(data, r.data[r.r:])
		return data
	}
	n1 := r.size - r.r
	copy(data[:n1], r.data[r.r:])
	copy(data[n1:], r.data[:r.w])
	return data
}

//可以用的长度
func (r *RingBuffer) Free() int {
	if r.IsFull() {
		return 0
	}
	if r.w < r.r {
		return r.r - r.w
	}
	return r.size - r.w + r.r
}
func (r *RingBuffer) IsFull() bool {
	/**
	 当读取时,一定是不满的,因为读取一定回去空出空间
	写入时,就不一定了,因为写入.能写满
	这个判断,读取时isEmpty一定为true导致条件不满足
	写入时.isEmpty为false,就需要判断读写,写指针是否相等,相等就是为满
	*/
	return !r.isEmpty && r.r == r.w
}
