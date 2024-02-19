package service

import "awesomeProject/opType"

type BindUserDealer struct {
}

func (dealer BindUserDealer) GetOperateType() int {
	return opType.BIND_USER
}

func (dealer BindUserDealer) DealOp(fromUid int, content []byte) {

}
