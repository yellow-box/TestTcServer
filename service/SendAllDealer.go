package service

import (
	"awesomeProject/entity"
	"awesomeProject/manager"
	"awesomeProject/opType"
	"encoding/json"
	"fmt"
)

type SendAllDealer struct {
}

func (dealer SendAllDealer) GetOperateType() int {
	return opType.SEND_MESSAGE_ALL
}

func (dealer SendAllDealer) DealOp(seq int64, fromUid int, content []byte) {
	if fromUid < 0 {
		fmt.Println("should bind user first uid:", fromUid)
		return
	}
	var roomMsg entity.RoomMsg
	err := json.Unmarshal(content, &roomMsg)
	if err != nil {
		fmt.Println("SendAllDealer unmarshal error:", err)
	} else {
		fmt.Printf("recieve from uid:%d, roomId:%d,conent:%s\n", fromUid, roomMsg.RoomId, roomMsg.Content)
		ucManager := manager.GetManager()
		chatRoom := GetChatRoomManager().GetChatRoom(roomMsg.RoomId)
		if chatRoom == nil {
			fmt.Println("chat Room not exist roomId:", roomMsg.RoomId)
			rsp := entity.Rsp{Code: opType.FAIL, Msg: "chat Room not exist"}
			rspByteArr, err := json.Marshal(rsp)
			if err != nil {
				return
			}
			ucManager.NotifyUserByte(fromUid, opType.CompoundContent(seq, dealer.GetOperateType(), rspByteArr))
			return
		}
		for _, uid := range chatRoom.uidS {
			if uid != fromUid {
				ucManager.NotifyUserByte(uid, opType.CompoundContent(seq, opType.PUSH_MESSAGE, content))
			} else {
				rsp := entity.Rsp{Code: opType.SUCCESS, Msg: ""}
				rspByteArr, err := json.Marshal(rsp)
				if err != nil {
					return
				}
				ucManager.NotifyUserByte(uid, opType.CompoundContent(seq, dealer.GetOperateType(), rspByteArr))
			}
		}
	}
}
