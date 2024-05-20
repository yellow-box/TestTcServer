package service

import (
	"awesomeProject/entity"
	"awesomeProject/manager"
	"awesomeProject/opType"
	"encoding/json"
	"fmt"
)

// 心跳包的处理
type HeartBeatDealer struct {
}

func (dealer HeartBeatDealer) GetOperateType() int {
	return opType.HEART_BEAT
}

func (dealer HeartBeatDealer) DealOp(seq int64, fromUid int, content []byte) {
	rsp, _ := json.Marshal(entity.Rsp{Code: opType.SUCCESS})
	fmt.Println("receive heat beat from uid:", fromUid)
	manager.GetManager().NotifyRecHeartBeat(fromUid)
	manager.GetManager().NotifyUserByte(fromUid, opType.CompoundContent(seq, opType.HEART_BEAT, rsp))
}
