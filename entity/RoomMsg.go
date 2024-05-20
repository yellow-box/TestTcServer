package entity

const System_Uid = 1

type RoomMsg struct {
	//房间id
	RoomId int `json:"room_id"`
	//消息内容
	Content string `json:"content"`
	//发送消息的id
	FromUid int `json:"from_uid"`
	//消息id
	MsgId string `json:"msg_id"`
}

type RoomInfo struct {
	//房间id
	RoomId int `json:"room_id"`
	//房间内的用户id
	UidS []int `json:"uid_s"`
}
