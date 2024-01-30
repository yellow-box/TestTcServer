package entry

import (
	"awesomeProject/manager"
	"bufio"
	"fmt"
	"net"
)

func startListen(conectInfo ConnectInfo) {
	addr := conectInfo.ip + ":" + conectInfo.port
	fmt.Printf("start listen :%s,name:%s", addr, conectInfo.name)
	listen, err := net.Listen(conectInfo.name, addr)
	if err != nil {
		fmt.Println("start listen occur error :", err)
		return
	}
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("\naccept clinet  occur error :", err)
			continue
		} else {
			fmt.Println("\naccept clinet  success: ", conn.RemoteAddr())
		}
		go acceptConnect(conn)
	}
}

func acceptConnect(conn net.Conn) {
	defer conn.Close()
	userConn := manager.UserConn{Conn: conn, Reader: bufio.NewReader(conn)}
	manager.GetManager().AppendUserConn(userConn)
	manager.GetManager().StartRead(userConn.Uid)
}

func Start() {
	cInfo := ConnectInfo{"10.30.10.114", "12345", "tcp"}
	startListen(cInfo)
}
