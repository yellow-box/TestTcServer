package manager

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
)

type User struct {
	Uid  int32
	Name string
}

type UserConn struct {
	Uid    int32
	Conn   net.Conn
	Reader *bufio.Reader
}

func (userConn UserConn) WriteString(content string) {
	userConn.WriteData([]byte(content))
}

func (userConn UserConn) WriteData(byteData []byte) {
	size := int32(len(byteData))
	resultSize := make([]byte, 4)
	var n int
	var err error
	binary.LittleEndian.PutUint32(resultSize, uint32(size))
	n, err = userConn.Conn.Write(resultSize)
	if err != nil {
		fmt.Println("write size to conn error:", err)
	} else {
		fmt.Println("write size success,num:", n)
	}
	n, err = userConn.Conn.Write(byteData)
	if err != nil {
		fmt.Println("write data to conn error:", err)
	} else {
		fmt.Println("write data success,num:", n)
	}
}

func (userConn UserConn) ReadData() []byte {
	//先读字节数
	dataBufSize := 1024
	buf := make([]byte, 128)
	var n int
	var err error
	var pendingSize int32
	n, err = userConn.Reader.Read(buf[:])
	if err != nil {
		fmt.Printf("userConn uid =%s ,read num error :%s\n", userConn.Uid, err)
	} else {
		fmt.Printf("userConn uid =%s, read num success \n")
	}
	pendingSize = getSizeFromByteArray(buf, n)
	readTurns := int(pendingSize) / dataBufSize
	if readTurns%dataBufSize != 0 {
		readTurns += 1
	}
	dataBuf := make([]byte, dataBufSize)
	for i := 0; i < readTurns; i++ {
		n, err = userConn.Reader.Read(dataBuf)
		if err != nil {
			fmt.Printf("userConn uid =%s ,read data error :%s\n", userConn.Uid, err)
		} else {
			fmt.Printf("userConn uid =%s, read data success \n")
		}
	}
	return dataBuf
}

func getSizeFromByteArray(rawContent []byte, endIndex int) int32 {
	return int32(binary.LittleEndian.Uint32(rawContent[:endIndex]))
}
