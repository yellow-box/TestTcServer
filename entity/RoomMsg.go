package entity

const System_Uid = 1

type RoomMsg struct {
	RoomId  int    `json:"room_id"`
	Content string `json:"content"`
	FromUid int    `json:"from_uid"`
	MsgId   string `json:"msg_id"`
}

type RoomInfo struct {
	RoomId int   `json:"room_id"`
	UidS   []int `json:"uid_s"`
}
