package service

import (
	"awesomeProject/manager"
	"awesomeProject/opType"
	"fmt"
)

type LogoutUserDealer struct {
}

func (dealer LogoutUserDealer) DealOp(fromUid int, content []byte) {
	if fromUid < 0 {
		fmt.Println("should bind user first uid:", fromUid)
		return
	}
	ucManager := manager.GetManager()
	ucManager.RemoveUserConnByUid(fromUid)
}

func (dealer LogoutUserDealer) GetOperateType() int {
	return opType.LOGOUT_UESR
}
