package manager

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"net"
)

type User struct {
	Uid  int
	Name string
}

type UserConn struct {
	Uid    int
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
	//tcp 粘包 问题，发送端 调用 api时分 多次调用，但在接受端 可能会合并成一次，来被接受。原因可能来自发送端，也可能来自接收端
}

func (userConn UserConn) readData(recCallback OnRecData) error {
	defer func(Conn net.Conn) {
		err := Conn.Close()
		if err != nil {
			fmt.Println("close Conn error uid=", userConn.Uid)
		}
	}(userConn.Conn)
	for {
		dataBytes, err := userConn.readDataInner()
		if err != nil {
			return err
		} else {
			recCallback.onRec(userConn.Uid, dataBytes)
		}
	}
}

func (userConn UserConn) readDataInner() ([]byte, error) {
	//先读字节数
	dataBufSize := 1024
	buf := make([]byte, 4)
	var n int
	var err error
	var pendingSize int
	n, err = userConn.Reader.Read(buf[:])
	if err != nil {
		fmt.Printf("\nuserConn uid =%d ,read num error :%s\n", userConn.Uid, err)
		return make([]byte, 0), err
	}
	pendingSize = getSizeFromByteArray(buf)
	fmt.Printf("\nuserConn uid =%d, read num success pendingSize=%d\n", userConn.Uid, pendingSize)
	readTurns := pendingSize / dataBufSize
	if pendingSize%dataBufSize != 0 {
		readTurns += 1
	}
	dataBuf := make([]byte, dataBufSize)
	result := make([]byte, pendingSize)
	for i := 0; i < readTurns; i++ {
		n, err = userConn.Reader.Read(dataBuf)
		if err != nil {
			fmt.Printf("userConn uid =%d ,read data error :%s\n", userConn.Uid, err)
			break
		} else {
			fmt.Printf("userConn uid =%d, read data success \n", userConn.Uid)
			startIndex := i * dataBufSize
			for j := 0; j < n; j++ {
				result[startIndex+j] = dataBuf[j]
			}
		}
	}
	return result, nil
}

func getSizeFromByteArray(rawContent []byte) int {
	return int(binary.LittleEndian.Uint32(rawContent))
}
