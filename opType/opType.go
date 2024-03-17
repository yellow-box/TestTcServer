package opType

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const oplen = 4
const seqlen = 8

// 规定最开始8个字节 表示 seq序号，接着4个字节表示 opType,后续字节表示具体内容
const BIND_USER = 1
const HEART_BEAT = 2
const SEND_MESSAGE_ALL = 3
const LOGOUT_UESR = 4
const PUSH_MESSAGE = 5
const CREATE_ROOM = 6
const JOIN_ROOM = 7

const SUCCESS = 200
const FAIL = 500

func divideThreePart(rawContent []byte) (error, int64, int, []byte) {
	if len(rawContent) < (oplen + seqlen) {
		return errors.New("content should start with 8 bytes as seq and then 4 bytes as optType"), 0, 0, make([]byte, 0)
	}
	headLen := seqlen + oplen
	seq := int64(binary.LittleEndian.Uint64(rawContent[:seqlen]))
	op := int(binary.LittleEndian.Uint32(rawContent[seqlen:headLen]))
	return nil, seq, op, rawContent[headLen:]
}

func CompoundContent(seq int64, opType int, content []byte) []byte {
	resultByte := make([]byte, len(content)+oplen+seqlen)
	binary.LittleEndian.PutUint64(resultByte, uint64(seq))
	binary.LittleEndian.PutUint32(resultByte[seqlen:seqlen+oplen], uint32(opType))
	for i := 0; i < len(content); i++ {
		resultByte[oplen+seqlen+i] = content[i]
	}
	return resultByte
}

type OpDealer interface {
	DealOp(seq int64, fromUid int, content []byte)
	GetOperateType() int
}

type RawMainDeal struct {
	OpDealMap map[int]OpDealer
}

func (mainDealer RawMainDeal) AddOpDealer(dealer OpDealer) {
	mainDealer.OpDealMap[dealer.GetOperateType()] = dealer
}

func (mainDealer RawMainDeal) Deal(fromUid int, rawContent []byte) {
	err, seq, operateType, content := divideThreePart(rawContent)
	if err != nil {
		fmt.Println("RawMainDeal deal error:", err)
		return
	}
	realDealer := mainDealer.OpDealMap[operateType]
	if realDealer != nil {
		realDealer.DealOp(seq, fromUid, content)
	} else {
		fmt.Println("no math opType:", operateType)
	}
}
