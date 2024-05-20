package service

import (
	"awesomeProject/entity"
	"awesomeProject/manager"
	"awesomeProject/opType"
	"encoding/json"
	"fmt"
)

// JoinRoomDealer 加入聊天室的处理
type JoinRoomDealer struct {
}

func (dealer JoinRoomDealer) GetOperateType() int {
	return opType.JOIN_ROOM
}

func (dealer JoinRoomDealer) DealOp(seq int64, fromUid int, content []byte) {
	var roomInfo entity.RoomInfo
	err := json.Unmarshal(content, &roomInfo)
	rsp := entity.Rsp{}
	if err != nil {
		fmt.Println("json Unmarshal error:", err)
		return
	}
	chatRoom := GetChatRoomManager().GetChatRoom(roomInfo.RoomId)
	ucManager := manager.GetManager()
	if chatRoom == nil {
		fmt.Println("join Room ,room not exist,roomId:", roomInfo.RoomId)
		rsp.Code = opType.FAIL
		rsp.Msg = "room not exist"
		rspArr, err := json.Marshal(rsp)
		if err != nil {
			return
		}
		ucManager.NotifyUserByte(fromUid, opType.CompoundContent(seq, dealer.GetOperateType(), rspArr))
	} else {
		chatRoom.addUid(fromUid)
		roomMsg := entity.RoomMsg{
			RoomId:  roomInfo.RoomId,
			FromUid: entity.System_Uid,
			Content: fmt.Sprint("uid:%d, join the room", fromUid),
		}
		roomMsgBytes, _ := json.Marshal(roomMsg)
		for _, uid := range chatRoom.uidS {
			ucManager.NotifyUserByte(uid, opType.CompoundContent(seq, opType.PUSH_MESSAGE, roomMsgBytes))
		}
	}
}
