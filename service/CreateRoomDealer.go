package service

import (
	"awesomeProject/entity"
	"awesomeProject/manager"
	"awesomeProject/opType"
	"encoding/json"
	"fmt"
)

// CreateRoomDealer 创建聊天室的处理
type CreateRoomDealer struct {
}

func (dealer CreateRoomDealer) GetOperateType() int {
	return opType.CREATE_ROOM
}

func (dealer CreateRoomDealer) DealOp(seq int64, fromUid int, content []byte) {
	var roomInfo entity.RoomInfo
	err := json.Unmarshal(content, &roomInfo)
	rsp := entity.Rsp{}
	if err != nil {
		fmt.Println("Unmarshal error :", err)
		return
	}
	chatRoom := GetChatRoomManager().GetChatRoom(roomInfo.RoomId)
	ucManager := manager.GetManager()
	if chatRoom != nil {
		fmt.Println("create Room fail, room has existed")
		rsp.Code = opType.FAIL
		rsp.Msg = "room has existed"
	} else {
		GetChatRoomManager().CreateRoom(roomInfo.RoomId)
		GetChatRoomManager().GetChatRoom(roomInfo.RoomId).addUid(fromUid)
		rsp.Code = opType.SUCCESS
		rsp.Msg = ""
	}
	rspArr, err := json.Marshal(rsp)
	if err != nil {
		return
	}
	ucManager.NotifyUserByte(fromUid, opType.CompoundContent(seq, dealer.GetOperateType(), rspArr))
}
