package opType

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const oplen = 4

// 规定前4个字节表示 opType,后续字节表示具体内容
const BIND_USER = 1
const SEND_MESSAGE_SINGEL = 2
const SEND_MESSAGE_ALL = 3
const LOGOUT_UESR = 4
const PUSH_MESSAGE = 5
const CREATE_ROOM = 6
const JOIN_ROOM = 7

func divideTwoPart(rawContent []byte) (error, int, []byte) {
	if len(rawContent) < oplen {
		return errors.New("content should start with 4 byte as opType"), 0, make([]byte, 0)
	}
	op := int(binary.LittleEndian.Uint32(rawContent[:oplen]))
	return nil, op, rawContent[4:]
}

func compoundContent(opType int, content []byte) []byte {
	resultByte := make([]byte, len(content)+oplen)
	binary.LittleEndian.PutUint32(resultByte, uint32(opType))
	for i := 0; i < len(content); i++ {
		resultByte[oplen+i] = content[i]
	}
	return resultByte
}

type OpDealer interface {
	DealOp(fromUid int, content []byte)
	GetOperateType() int
}

type RawMainDeal struct {
	OpDealMap map[int]OpDealer
}

func (mainDealer RawMainDeal) AddOpDealer(dealer OpDealer) {
	mainDealer.OpDealMap[dealer.GetOperateType()] = dealer
}

func (mainDealer RawMainDeal) Deal(fromUid int, rawContent []byte) {
	err, operateType, content := divideTwoPart(rawContent)
	if err != nil {
		fmt.Println("RawMainDeal deal error:", err)
		return
	}
	realDealer := mainDealer.OpDealMap[operateType]
	if realDealer != nil {
		realDealer.DealOp(fromUid, content)
	} else {
		fmt.Println("no math opType:", operateType)
	}
}
