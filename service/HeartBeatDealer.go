package service

import (
	"awesomeProject/manager"
	"awesomeProject/opType"
)

type HeartBeatDealer struct {
}

func (dealer HeartBeatDealer) GetOperateType() int {
	return opType.HEART_BEAT
}

func (dealer HeartBeatDealer) DealOp(seq int64, fromUid int, content []byte) {
	manager.GetManager().NotifyUserByte(fromUid, opType.CompoundContent(seq, opType.HEART_BEAT, make([]byte, 0)))
}
