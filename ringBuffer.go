package ringbuffer

import (
	"bytes"
	"errors"
)

var (
	errFull  = errors.New("已经最大空间")
	errEmpty = errors.New("缓存为空")
)

/**
不是使用虚拟指针,
因为一般的分包都是 payload分隔符 或者长度+payload 或者 分隔符payload分隔符
所以通过2个方法,在不修改读取指针的情况下查询结果
新增成为动态扩容的环形缓存,新增一个maxPacketSize 作为一个环形缓存的最大缓存大小防止错误的分包
导致的内存爆炸..
*/
type RingBuffer struct {
	data          []byte
	r             int //下一个读取的位置
	w             int //下一个写入的位置
	size          int
	isEmpty       bool
	maxPacketSize int //默认为1024000字节的长度
	autoSize      int //每次扩大的容量大小
}

func NewRingBuffer(length int) *RingBuffer {
	return &RingBuffer{
		data:          make([]byte, length),
		r:             0,
		w:             0,
		size:          length,
		isEmpty:       true,
		maxPacketSize: 1024 * 1000, //默认为1024000字节的长度
		autoSize:      1024,        //每次扩大的容量大小
	}
}
func (r *RingBuffer) SetMaxSize(size int) {
	r.maxPacketSize = size
}
func (r *RingBuffer) SetAutoSize(size int) {
	r.autoSize = size
}

//读取一段,但是读取指针不会该改变
func (r *RingBuffer) VirtualRead(target []byte) (int, error) {
	return r.read(target)
}

//读取到分割符号,可以跳过的
func (r *RingBuffer) Index(sep []byte, step int) int {
	if r.w > r.r {
		return bytes.Index(r.data[r.r+step:r.w], sep) + step
	}
	i := bytes.Index(r.data[r.r+step:], sep)
	if i == -1 {
		i = bytes.Index(r.data[:r.w], sep) + r.size - r.r
	}
	if i == -1 {
		return -1
	}
	return i + step
}

func (r *RingBuffer) read(target []byte) (int, error) {
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
	} else {
		c1 := r.size - r.r
		if c1 >= n {
			copy(target, r.data[r.r:])
		} else {
			copy(target, r.data[r.r:])
			copy(target[c1:], r.data[:n-r.size+r.r])
		}
	}
	return n, nil
}
func (r *RingBuffer) Read(target []byte) (int, error) {
	n, err := r.read(target)
	if err != nil {
		return 0, err
	}
	r.r = (r.r + n) % r.size
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
	free := r.Free()
	n := len(data)
	if n > free {
		//这个时候也需要扩大缓存
		if r.size >= r.maxPacketSize {
			return 0, errFull
		}
		r.makeSpace(n - free)
	}
	if r.r > r.w {
		copy(r.data[r.w:], data)
		r.w += n
	} else {
		c1 := r.size - r.w
		if c1 >= n {
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
	for r.IsFull() {
		if r.size >= r.maxPacketSize {
			return errFull
		}
		r.size += r.autoSize
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

func (r *RingBuffer) makeSpace(n int) {
	//我草,,可以不去计算咋个做
	//直接移动到新的数组中,也对反正长度新增了,也需要构筑新的数组
	newSize := r.size + n
	oldLength := r.Length()
	newData := make([]byte, newSize)
	_, _ = r.Read(newData)
	r.r = 0
	r.w = oldLength
	r.data = newData
	r.size = newSize
}
