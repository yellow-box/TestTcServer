package entry

import (
	"awesomeProject/manager"
	"awesomeProject/opType"
	"awesomeProject/service"
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
	userConn := manager.UserConn{Conn: conn, Reader: bufio.NewReader(conn), Uid: -1}
	err := userConn.BindUser()
	if err != nil {
		fmt.Println("bindUser error:", err)
		return
	}
	manager.GetManager().StartRead(userConn.Uid)
}

func Start() {
	cInfo := ConnectInfo{"10.30.10.114", "12345", "tcp"}
	service.ChatService{MainDealer: opType.RawMainDeal{OpDealMap: make(map[int]opType.OpDealer)}}.Init()
	startListen(cInfo)
}
