package manager

import (
	"sync"
)

var once sync.Once
var manager UserConnManager

type UserConnManager struct {
	userConnMap map[int32]UserConn
}

func GetManager() {
	once.Do(func() {
		manager = UserConnManager{make(map[int32]UserConn)}
	})
}

func (ucManager UserConnManager) appendUserConn(uConn UserConn) {
	ucMap := ucManager.userConnMap
	ucMap[uConn.Uid] = uConn
}

func (ucManager UserConnManager) removeUserConn(uConn UserConn) {
	ucMap := ucManager.userConnMap
	delete(ucMap, uConn.Uid)
}
