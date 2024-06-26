package manager

import (
	"fmt"
	"sync"
	"time"
)

var once sync.Once
var manager UserConnManager
var recCallbackList []OnRecData

// UserConnManager 用户套接字管理
type UserConnManager struct {
	userConnMap map[int]UserConn
	rwLock      sync.RWMutex
}

// OnRecData 收到用户发送消息后的回调接口
type OnRecData interface {
	OnRec(userConn *UserConn, data []byte)
}

// AddRecCallback 存储收到用户套接字后的 回调
func AddRecCallback(onRecData OnRecData) {
	recCallbackList = append(recCallbackList, onRecData)
}

func (ucManager *UserConnManager) OnRec(uConn *UserConn, data []byte) {
	fmt.Printf("recv data uid:%d,data:%s\n", uConn.Uid, string(data[:]))
	for _, recCallback := range recCallbackList {
		recCallback.OnRec(uConn, data)
	}
}

func GetManager() *UserConnManager {
	once.Do(func() {
		manager = UserConnManager{userConnMap: make(map[int]UserConn)}
	})
	return &manager
}

// AppendUserConn 新增一个用户套接字
func (ucManager *UserConnManager) AppendUserConn(uConn *UserConn) {
	defer ucManager.rwLock.Unlock()
	ucManager.rwLock.Lock()
	ucManager.userConnMap[uConn.Uid] = *uConn
}

func (ucManager *UserConnManager) DeleteUserConn(uConn *UserConn) {
	defer ucManager.rwLock.Unlock()
	ucManager.rwLock.Lock()
	//等于 -1，说明还没bindUser成功
	if uConn.Uid != -1 {
		delete(ucManager.userConnMap, uConn.Uid)
	} else {
		for k, userConn := range ucManager.userConnMap {
			if &userConn == uConn {
				delete(ucManager.userConnMap, k)
			}
		}
	}
	ucManager.userConnMap[uConn.Uid] = *uConn
}

func (ucManager *UserConnManager) StartRead(uid int) {
	uConn := ucManager.userConnMap[uid]
	if uConn.Conn == nil {
		fmt.Println("no conn match uid:", uid)
		return
	}
	err := uConn.readData(&manager)
	if err != nil {
		fmt.Printf("read Data error uid=%d,error=%s\n", uid, err)
		ucManager.RemoveUserConn(uConn)
	}
}

// StartReadConn 读取用户套接字内容
func (ucManager *UserConnManager) StartReadConn(uConn UserConn) {
	err := uConn.readData(&manager)
	if err != nil {
		fmt.Printf("read Data error uid=%d,error=%s\n", uConn.Uid, err)
		ucManager.RemoveUserConn(uConn)
	}
}

func (ucManager *UserConnManager) RemoveUserConn(uConn UserConn) {
	defer ucManager.rwLock.Unlock()
	ucMap := ucManager.userConnMap
	ucManager.rwLock.Lock()
	delete(ucMap, uConn.Uid)
}

func (ucManager *UserConnManager) RemoveUserConnByUid(uid int) {
	defer ucManager.rwLock.Unlock()
	ucManager.rwLock.Lock()
	uConn := ucManager.userConnMap[uid]
	if uConn.Conn != nil {
		err := uConn.Conn.Close()
		if err != nil {
			delete(ucManager.userConnMap, uid)
			fmt.Println("success remove uConn and close conn")
		}
	}
}

// NotifyAll 为所有用户发送消息
func (ucManager *UserConnManager) NotifyAll(content string) {
	defer ucManager.rwLock.RUnlock()
	ucManager.rwLock.RLock()
	for _, uConn := range ucManager.userConnMap {
		uConn.WriteString(content)
	}
}

// NotifyUser 为特定用户发送消息
func (ucManager *UserConnManager) NotifyUser(uid int, content string) {
	defer ucManager.rwLock.RUnlock()
	ucManager.rwLock.RLock()
	uConn := ucManager.userConnMap[(uid)]
	uConn.WriteString(content)
}

// NotifyUserByte 为特定用户发送消息
func (ucManager *UserConnManager) NotifyUserByte(uid int, content []byte) {
	defer ucManager.rwLock.RUnlock()
	ucManager.rwLock.RLock()
	conn, ok := ucManager.userConnMap[(uid)]
	if ok {
		conn.WriteData(content)
	}
}

// NotifyRecHeartBeat 为用户发送心跳包回包
func (ucManager *UserConnManager) NotifyRecHeartBeat(uid int) {
	defer ucManager.rwLock.RUnlock()
	ucManager.rwLock.RLock()
	conn, ok := ucManager.userConnMap[(uid)]
	if ok {
		conn.LastHeartBeatRecTime = time.Now().UnixNano()
	}
}
