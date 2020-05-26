package ringbuffer

import "errors"

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
	if r.r == r.w {
		r.isEmpty = false
	}
	return n, nil
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
	if r.w == r.r {
		r.isEmpty = false
	}
	return n, nil
}

//保存的数据长度
func (r *RingBuffer) Length() int {
	if r.r <= r.w {
		return r.w - r.r
	}
	return r.size - r.r + r.w
}

//可以用的长度
func (r *RingBuffer) Free() int {
	if r.w < r.r {
		return r.r - r.w
	}
	return r.size - r.w + r.r
}
func (r *RingBuffer) IsFull() bool {
	return !r.isEmpty && r.r == r.w
}
