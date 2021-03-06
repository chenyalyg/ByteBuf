package bytebuf

import (
	"container/list"
	"encoding/binary"
	"errors"
	"math"
	"sync"
)

type CopyType int32

const (
	DeepCopy      CopyType = 0 //创建新的内存存放数据
	Pointer      CopyType = 1 //
)

type ByteBuffer struct {
	cache    *list.List
	len      int
	mu       sync.RWMutex
	copy_type CopyType
}

//copy_type 是否进行深度拷贝
func New(copy_type CopyType) *ByteBuffer {
	b_buf:= &ByteBuffer{}
	b_buf.cache=list.New()
	b_buf.copy_type=copy_type
	return b_buf
}



func (self *ByteBuffer)Len() int {
	return self.len
}

func (self *ByteBuffer) ReadBytes(data []byte) int {
	self.mu.RLock()
	defer self.mu.RUnlock()

	if self.len==0 {
		return 0
	}

	var index=0
	for f:=self.cache.Front();f!=nil;{
		c_buf:=f.Value.([]byte)
		c_node:=f
		f=f.Next()

		n:=copy(data[index:],c_buf)
		index+=n
		if n== len(c_buf) {
			self.cache.Remove(c_node)
		}else{
			c_buf=c_buf[n:]
			self.cache.InsertBefore(c_buf,c_node)
			self.cache.Remove(c_node)
		}
		if index>= len(data) {
			break
		}
	}

	self.len-=index
	return index
}

//预读取字节
func (self *ByteBuffer) PrvReadBytes(data []byte) int {
	self.mu.RLock()
	defer self.mu.RUnlock()

	if self.len==0 {
		return 0
	}

	var index=0
	for f:=self.cache.Front();f!=nil;{
		c_buf:=f.Value.([]byte)
		f=f.Next()

		n:=copy(data[index:],c_buf)
		index+=n
		if index>= len(data) {
			break
		}
	}

	return index
}

func (self*ByteBuffer)check_len(size int) error {
	if self.len==0 {
		return errors.New("buffer is empty!")
	}
	if self.len<size {
		return errors.New("buffer is low!")
	}
	return nil
}

func (self *ByteBuffer) ReadByte() (byte,error) {
	err:=self.check_len(1)
	if err!=nil {
		return 0,err
	}
	buf:=make([]byte,1)
	self.ReadBytes(buf)
	return buf[0],nil
}

func (self *ByteBuffer) PrvReadByte() (byte,error) {
	err:=self.check_len(1)
	if err!=nil {
		return 0,err
	}
	buf:=make([]byte,1)
	self.PrvReadBytes(buf)
	return buf[0],nil
}

func (self *ByteBuffer) ReadInt16(order binary.ByteOrder) (int16,error) {
	err:=self.check_len(2)
	if err!=nil {
		return 0,err
	}
	buf:=make([]byte,2)
	self.ReadBytes(buf)
	return int16(order.Uint16(buf)),nil
}

func (self *ByteBuffer) PrvReadInt16(order binary.ByteOrder) (int16,error) {
	err:=self.check_len(2)
	if err!=nil {
		return 0,err
	}
	buf:=make([]byte,2)
	self.PrvReadBytes(buf)
	return int16(order.Uint16(buf)),nil
}

func (self *ByteBuffer) ReadInt32(order binary.ByteOrder) (int32,error) {
	err:=self.check_len(4)
	if err!=nil {
		return 0,err
	}
	buf:=make([]byte,4)
	self.ReadBytes(buf)
	return int32(order.Uint32(buf)),nil
}

func (self *ByteBuffer) PrvReadInt32(order binary.ByteOrder) (int32,error) {
	err:=self.check_len(4)
	if err!=nil {
		return 0,err
	}
	buf:=make([]byte,4)
	self.PrvReadBytes(buf)
	return int32(order.Uint32(buf)),nil
}

func (self *ByteBuffer) ReadInt64(order binary.ByteOrder) (int64,error) {
	err:=self.check_len(8)
	if err!=nil {
		return 0,err
	}
	buf:=make([]byte,8)
	self.ReadBytes(buf)
	return int64(order.Uint64(buf)),nil
}

func (self *ByteBuffer) PrvReadInt64(order binary.ByteOrder) (int64,error) {
	err:=self.check_len(8)
	if err!=nil {
		return 0,err
	}
	buf:=make([]byte,8)
	self.PrvReadBytes(buf)
	return int64(order.Uint64(buf)),nil
}

func (self *ByteBuffer) ReadFloat32(order binary.ByteOrder) (float32,error) {
	err:=self.check_len(4)
	if err!=nil {
		return 0,err
	}
	buf:=make([]byte,4)
	self.ReadBytes(buf)

	bits := order.Uint32(buf)
	return math.Float32frombits(bits),nil
}

func (self *ByteBuffer) PrvReadFloat32(order binary.ByteOrder) (float32,error) {
	err:=self.check_len(4)
	if err!=nil {
		return 0,err
	}
	buf:=make([]byte,4)
	self.PrvReadBytes(buf)

	bits := order.Uint32(buf)
	return math.Float32frombits(bits),nil
}

func (self *ByteBuffer) ReadFloat64(order binary.ByteOrder) (float32,error) {
	err:=self.check_len(8)
	if err!=nil {
		return 0,err
	}
	buf:=make([]byte,8)
	self.ReadBytes(buf)

	bits := order.Uint32(buf)
	return math.Float32frombits(bits),nil
}

func (self *ByteBuffer) PrvReadFloat64(order binary.ByteOrder) (float32,error) {
	err:=self.check_len(8)
	if err!=nil {
		return 0,err
	}
	buf:=make([]byte,8)
	self.PrvReadBytes(buf)

	bits := order.Uint32(buf)
	return math.Float32frombits(bits),nil
}

func (self *ByteBuffer)WriteBytes(data []byte) {
	var copy_buf=data
	if self.copy_type==DeepCopy {
		copy_buf= make([]byte, len(data))
		copy(copy_buf,data)
	}
	self.mu.Lock()
	self.cache.PushBack(copy_buf)
	self.len+= len(copy_buf)
	self.mu.Unlock()
}

func (self *ByteBuffer)WriteByte(b byte)  {
	buf:=make([]byte,1)
	buf[0]=b
	self.WriteBytes(buf)
}

func (self *ByteBuffer)WriteInt16(num int16,order binary.ByteOrder)  {
	buf:=make([]byte,2)
	order.PutUint16(buf,uint16(num))
	self.WriteBytes(buf)
}

func (self *ByteBuffer)WriteInt32(num int32,order binary.ByteOrder)  {
	buf:=make([]byte,4)
	order.PutUint32(buf,uint32(num))
	self.WriteBytes(buf)
}

func (self *ByteBuffer)WriteInt64(num int64,order binary.ByteOrder)  {
	buf:=make([]byte,8)
	order.PutUint64(buf,uint64(num))
	self.WriteBytes(buf)
}

func (self *ByteBuffer)WriteFloat32(num float32,order binary.ByteOrder)  {
	bits :=math.Float32bits(num)
	buf:=make([]byte,4)
	order.PutUint32(buf,bits)
	self.WriteBytes(buf)
}

func (self *ByteBuffer)WriteFloat64(num float64,order binary.ByteOrder)  {
	bits :=math.Float64bits(num)
	buf:=make([]byte,8)
	order.PutUint64(buf,bits)
	self.WriteBytes(buf)
}