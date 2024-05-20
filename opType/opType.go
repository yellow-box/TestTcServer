package opType

import (
	"awesomeProject/manager"
	"encoding/binary"
	"errors"
	"fmt"
)

const oplen = 4
const seqlen = 8

// 规定最开始8个字节 表示 seq序号，接着4个字节表示 opType,后续字节表示具体内容

// BIND_USER 给套接字绑定用户
const BIND_USER = 1

// HEART_BEAT 发送心跳包
const HEART_BEAT = 2

// SEND_MESSAGE_ALL 发送消息
const SEND_MESSAGE_ALL = 3

// LOGOUT_UESR 退出登录
const LOGOUT_UESR = 4

// PUSH_MESSAGE 发送push消息
const PUSH_MESSAGE = 5

// CREATE_ROOM 创建一个聊天室
const CREATE_ROOM = 6

// JOIN_ROOM 加入一个聊天室
const JOIN_ROOM = 7

// SUCCESS 在回包中添加，表示请求成功
const SUCCESS = 200

// FAIL 在回包中添加，表示请求失败
const FAIL = 500

// HeartBeatInterval 心跳包检测间隔
const HeartBeatInterval = 1000

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

type BindUserAction interface {
	BindUser(conn *manager.UserConn, data []byte)
}

type RawMainDeal struct {
	OpDealMap map[int]OpDealer
}

func (mainDealer RawMainDeal) AddOpDealer(dealer OpDealer) {
	mainDealer.OpDealMap[dealer.GetOperateType()] = dealer
}

func (mainDealer RawMainDeal) Deal(userConn *manager.UserConn, bindUserAction BindUserAction, rawContent []byte) {
	fromUid := userConn.Uid
	err, seq, operateType, content := divideThreePart(rawContent)
	if err != nil {
		fmt.Println("RawMainDeal deal error:", err)
		return
	}
	if operateType == BIND_USER {
		bindUserAction.BindUser(userConn, content)
		return
	}
	realDealer := mainDealer.OpDealMap[operateType]
	if realDealer != nil {
		realDealer.DealOp(seq, fromUid, content)
	} else {
		fmt.Println("no math opType:", operateType)
	}
}
