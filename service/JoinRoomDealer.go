package service

import (
	"awesomeProject/entity"
	"awesomeProject/opType"
	"encoding/json"
	"fmt"
)

type JoinRoomDealer struct {
}

func (dealer JoinRoomDealer) GetOperateType() int {
	return opType.JOIN_ROOM
}

func (dealer JoinRoomDealer) DealOp(fromUid int, content []byte) {
	var roomInfo entity.RoomInfo
	err := json.Unmarshal(content, &roomInfo)
	if err != nil {
		fmt.Println("json Unmarshal error:", err)
		return
	}
	chatRoom := GetChatRoomManager().GetChatRoom(roomInfo.RoomId)
	if chatRoom == nil {
		fmt.Println("join Room ,room not exist,roomId:", roomInfo.RoomId)
	} else {
		chatRoom.addUid(fromUid)
	}
}
