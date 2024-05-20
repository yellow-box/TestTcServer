package service

import "sync"

var once sync.Once
var chatRoomManager ChatRoomManager

// 聊天室管理的抽象
type ChatRoomManager struct {
	chatRoomMap map[int]*ChatRoom
	rwLock      sync.RWMutex
}

func GetChatRoomManager() *ChatRoomManager {
	once.Do(func() {
		chatRoomManager = ChatRoomManager{chatRoomMap: make(map[int]*ChatRoom), rwLock: sync.RWMutex{}}
		chatRoomManager.InitTestRoom()
	})
	return &chatRoomManager
}

func (manager *ChatRoomManager) GetChatRoom(roomId int) *ChatRoom {
	defer manager.rwLock.RUnlock()
	manager.rwLock.RLock()
	return manager.chatRoomMap[roomId]
}

func (manager *ChatRoomManager) CreateRoom(tRoomId int) {
	defer manager.rwLock.Unlock()
	manager.rwLock.Lock()
	manager.chatRoomMap[tRoomId] = &ChatRoom{roomId: tRoomId, uidS: make([]int, 0)}

}

func (manager *ChatRoomManager) InitTestRoom() {
	testRoomId := 1
	manager.CreateRoom(testRoomId)
	testRoom := manager.GetChatRoom(testRoomId)
	testRoom.addUid(3)
	testRoom.addUid(4)
	testRoom.addUid(5)
}
