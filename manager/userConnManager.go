package manager

import (
	"fmt"
	"sync"
)

var once sync.Once
var manager UserConnManager
var recCallback OnRecData

type UserConnManager struct {
	userConnMap map[int]UserConn
	recCallback OnRecData
}

type OnRecData interface {
	onRec(uid int, data []byte)
}

type CallDealer struct {
}

func (deal CallDealer) onRec(uid int, data []byte) {
	fmt.Printf("recv data uid:%d,data:%s\n", uid, string(data[:]))
	GetManager().NotifyUserByte(uid, data)
}

func GetManager() UserConnManager {
	once.Do(func() {
		manager = UserConnManager{userConnMap: make(map[int]UserConn), recCallback: CallDealer{}}
	})
	return manager
}

func (ucManager UserConnManager) AppendUserConn(uConn UserConn) {
	ucMap := ucManager.userConnMap
	ucMap[uConn.Uid] = uConn
}

func (ucManager UserConnManager) StartRead(uid int) {
	uconn := ucManager.userConnMap[uid]
	err := uconn.readData(manager.recCallback)
	if err != nil {
		fmt.Printf("read Data error uid=%d,error=%s\n", uid, err)
		ucManager.RemoveUserConn(uconn)
	}

}

func (ucManager UserConnManager) RemoveUserConn(uConn UserConn) {
	ucMap := ucManager.userConnMap
	delete(ucMap, uConn.Uid)
}

func (ucManager UserConnManager) NotifyAll(content string) {
	for _, userconn := range ucManager.userConnMap {
		userconn.WriteString(content)
	}
}

func (ucManager UserConnManager) NotifyUser(uid int, content string) {
	conn := ucManager.userConnMap[(uid)]
	conn.WriteString(content)
}

func (ucManager UserConnManager) NotifyUserByte(uid int, content []byte) {
	conn := ucManager.userConnMap[(uid)]
	conn.WriteData(content)
}
