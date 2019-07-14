package bytebuf

import (
	"encoding/binary"
	"fmt"
	"time"
)

func init() {
	test1()
	test2()
}

func test1()  {
	start_time:=time.Now().Unix()
	b_buf:=New(false)

	for i:=0;i<9999999;i++{
		buf:=make([]byte,1500)
		b_buf.WriteBytes(buf)
		//b_buf.WriteInt16(12345,binary.BigEndian)
	}

	fmt.Println(b_buf.ReadInt16(binary.BigEndian))
	fmt.Println(b_buf.ReadInt16(binary.BigEndian))


	finish_time:=time.Now().Unix()

	fmt.Println("use time:",finish_time-start_time)
}


func test2()  {
	start_time:=time.Now().Unix()
	b_buf:=New(false)
	b_buf.WriteFloat32(12345.22,binary.BigEndian)
	for i:=0;i<999999;i++{
		b_buf.WriteInt16(12345,binary.BigEndian)
		//b_buf.WriteInt16(12345,binary.BigEndian)
	}
	b_buf.WriteInt32(12345,binary.BigEndian)
	b_buf.WriteByte(123)

	num,_:=b_buf.ReadFloat32(binary.BigEndian)
	fmt.Println("num :",num)
	buf:=make([]byte,999999*3)
	n:=b_buf.ReadBytes(buf)
	fmt.Println("read size:",n)

	finish_time:=time.Now().Unix()

	fmt.Println("use time:",finish_time-start_time)
}