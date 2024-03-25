package manager

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"time"
)

type User struct {
	Uid  int
	Name string
}

type UserConn struct {
	Uid                  int
	Conn                 net.Conn
	Reader               *bufio.Reader
	LastHeartBeatRecTime int64
}

func (userConn *UserConn) WriteString(content string) {
	userConn.WriteData([]byte(content))
}

func SetDefaultUserConn(Uid int, Conn net.Conn, Reader *bufio.Reader) UserConn {
	return UserConn{Uid: Uid, Conn: Conn, Reader: Reader, LastHeartBeatRecTime: time.Now().UnixNano()}
}

func (userConn *UserConn) WriteData(byteData []byte) {
	if userConn.Conn == nil {
		fmt.Println("uConn is not exist,uid:", userConn.Uid)
		return
	}
	size := int32(len(byteData))
	resultSize := make([]byte, 4)
	var n int
	var err error
	binary.LittleEndian.PutUint32(resultSize, uint32(size))
	//tcpConn := userConn.Conn.(*net.TCPConn)
	//err = tcpConn.SetNoDelay(true)
	//if err != nil {
	//	return
	//}
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

func (userConn *UserConn) readData(recCallback OnRecData) error {
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
			recCallback.OnRec(userConn, dataBytes)
		}
	}
}

func (userConn *UserConn) readDataInner() ([]byte, error) {
	//先读字节数
	dataBufSize := 1024
	buf := make([]byte, 4)
	var n int
	var err error
	var pendingSize int
	n, err = io.ReadFull(userConn.Reader, buf[:])
	if err != nil {
		fmt.Printf("\nuserConn uid =%d ,read num error :%s\n", userConn.Uid, err)
		return make([]byte, 0), err
	}
	pendingSize = getSizeFromByteArray(buf)
	fmt.Printf("\nuserConn uid =%d, read num success pendingSize=%d\n", userConn.Uid, pendingSize)
	curReadSize := 0
	dataBuf := make([]byte, dataBufSize)
	result := make([]byte, pendingSize)
	for curReadSize < pendingSize {
		if (pendingSize - curReadSize) < dataBufSize {
			dataBuf = dataBuf[:pendingSize-curReadSize]
		}
		n, err = io.ReadFull(userConn.Reader, dataBuf)
		if err != nil {
			fmt.Printf("userConn uid =%d ,read data error :%s\n", userConn.Uid, err)
			return result, err
		} else {
			fmt.Printf("userConn uid =%d, read data success \n", userConn.Uid)
			for j := 0; j < n; j++ {
				result[curReadSize+j] = dataBuf[j]
			}
			curReadSize += n
		}
	}
	return result, nil
}

func getSizeFromByteArray(rawContent []byte) int {
	return int(binary.LittleEndian.Uint32(rawContent))
}

func (userConn *UserConn) BindUser() error {
	rawContent, err := userConn.readDataInner()
	if err != nil {
		return err
	}
	uid := binary.LittleEndian.Uint32(rawContent)
	if uid < 1 {
		return errors.New(" uid should greater than 0,uid")
	}
	userConn.Uid = int(uid)
	GetManager().AppendUserConn(userConn)
	return nil
}
