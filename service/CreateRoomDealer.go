package service

import (
	"awesomeProject/entity"
	"awesomeProject/opType"
	"encoding/json"
	"fmt"
)

type CreateRoomDealer struct {
}

func (dealer CreateRoomDealer) GetOperateType() int {
	return opType.CREATE_ROOM
}

func (dealer CreateRoomDealer) DealOp(fromUid int, content []byte) {
	var roomInfo entity.RoomInfo
	err := json.Unmarshal(content, &roomInfo)
	if err != nil {
		fmt.Println("Unmarshal error :", err)
		return
	}
	chatRoom := GetChatRoomManager().GetChatRoom(roomInfo.RoomId)
	if chatRoom != nil {
		fmt.Println("create Room fail, room has existed")
		return
	}
	GetChatRoomManager().CreateRoom(roomInfo.RoomId)
	GetChatRoomManager().GetChatRoom(roomInfo.RoomId).addUid(fromUid)
}
