# ByteBuf
Speed byte cache

## Installing
```bash
go get -u github.com/chenyalyg/ByteBuf
```

## Usage

```go
package main

import (
  "fmt"
  "log"
  
  "github.com/chenyalyg/ByteBuf"
)

func main() {
  buf := bytebuf.New(false)
  buf.WriteInt16(12345,binary.BigEndian)
  buf.WriteInt32(12345,binary.BigEndian)
  buf.WriteFloat32(12345.22,binary.BigEndian)
  buf.WriteFloat64(12345.33,binary.BigEndian)
  b:=make([]byte,10)
  buf.WriteBytes(b)
  
  num1,err:=buf.ReadInt16(binary.BigEndian)
  if err != nil {
	log.Fatal(err)
  }
  fmt.Println("num1 :",num1)
  
  num2,err:=buf.ReadInt32(binary.BigEndian)
  if err != nil {
	log.Fatal(err)
  }
  fmt.Println("num2 :",num2)
  
  num3,err:=buf.ReadFloat32(binary.BigEndian)
  if err != nil {
	log.Fatal(err)
  }
  fmt.Println("num3 :",num3)
  
  num4,err:=buf.ReadFloat64(binary.BigEndian)
  if err != nil {
	log.Fatal(err)
  }
  fmt.Println("num4 :",num4)
  
  b2:=make([]byte,buf.Len())
  n:=buf.ReadBytes(b2)
  fmt.Println("read size:",n)

}
```
