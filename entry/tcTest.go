package entry

import (
	"awesomeProject/manager"
	"awesomeProject/opType"
	"awesomeProject/service"
	"bufio"
	"fmt"
	"net"
	"time"
)

func startListen(connectInfo ConnectInfo) {
	addr := connectInfo.ip + ":" + connectInfo.port
	fmt.Printf("start listen :%s,name:%s", addr, connectInfo.name)
	listen, err := net.Listen(connectInfo.name, addr)
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
	userConn := manager.SetDefaultUserConn(-1, conn, bufio.NewReader(conn))
	heatBeatListen(&userConn)
	manager.GetManager().StartReadConn(userConn)
}

func heatBeatListen(conn *manager.UserConn) {
	// 定义定时器，每隔1秒执行一次任务
	ticker := time.NewTicker(opType.HeartBeatInterval * time.Millisecond)
	fmt.Println("start heatBeatListen")
	go func(conn *manager.UserConn) {
		defer func(conn *manager.UserConn) {
			if conn == nil {
				return
			}
			conn.Conn.Close()
		}(conn)
		for {
			//阻塞，直到定时器触发并且时间值被发送到 ticker.C 通道
			<-ticker.C
			cur := time.Now().UnixNano() //当前时间戳  s
			//假设 RRT/2 = 500s
			//超时未收到 心跳包
			//fmt.Println("heartbeat check")
			if (cur - conn.LastHeartBeatRecTime) > ((opType.HeartBeatInterval + 500) * int64(time.Millisecond)) {
				fmt.Println("long time not receive heat beat, ready to delete this conn,ConnUid:", conn.Uid)
				manager.GetManager().DeleteUserConn(conn)
				break
			} else {
				//fmt.Println("heartbeat check cur =", cur)
				conn.LastHeartBeatRecTime = cur
			}
		}
	}(conn)
}

func Start() {
	cInfo := ConnectInfo{"", "12345", "tcp"}
	service.ChatService{MainDealer: opType.RawMainDeal{OpDealMap: make(map[int]opType.OpDealer)}}.Init()
	startListen(cInfo)
}
