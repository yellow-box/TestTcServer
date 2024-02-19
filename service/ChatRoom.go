package service

type ChatRoom struct {
	roomId int
	uidS   []int
}

func (room *ChatRoom) addUid(uid int) {
	if !room.userInRoom(uid) {
		room.uidS = append(room.uidS, uid)
	}
}

func (room *ChatRoom) userInRoom(tUid int) bool {
	for _, uid := range room.uidS {
		if tUid == uid {
			return true
		}
	}
	return false
}

func (room *ChatRoom) deleteUid(tUid int) {
	i := 0
	for _, uid := range room.uidS {
		if uid != tUid {
			room.uidS[i] = uid
			i++
		}
	}
	room.uidS = room.uidS[:i]
}
