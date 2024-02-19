package entity

type RoomMsg struct {
	RoomId  int    `json:"room_id"`
	Content string `json:"content"`
}

type RoomInfo struct {
	RoomId int   `json:"room_id"`
	UidS   []int `json:"uid_s"`
}
