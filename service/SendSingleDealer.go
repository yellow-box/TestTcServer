package service

import (
	"awesomeProject/manager"
	"awesomeProject/opType"
	"fmt"
)

type SendSingleDealer struct {
}

func (dealer SendSingleDealer) GetOperateType() int {
	return opType.SEND_MESSAGE_SINGEL
}

func (dealer SendSingleDealer) DealOp(fromUid int, content []byte) {
	if fromUid < 0 {
		fmt.Println("should bind user first uid:", fromUid)
		return
	}
	ucManager := manager.GetManager()
	ucManager.NotifyUser(fromUid, string(content[:]))
}
