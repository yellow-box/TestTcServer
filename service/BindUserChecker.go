package service

import (
	"awesomeProject/manager"
	"encoding/json"
	"fmt"
)

// BinUserData 绑定用户信息
type BinUserData struct {
	Uid int `json:"uid"`
}

type BindUserCheck struct {
}

func (b BindUserCheck) BindUser(userConn *manager.UserConn, data []byte) {
	var bindUserData BinUserData
	json.Unmarshal(data, &bindUserData)
	if bindUserData.Uid < 1 {
		fmt.Printf("bindUser error,uid:%d, it must larger than 0\n", bindUserData.Uid)
	} else {
		userConn.Uid = bindUserData.Uid
		manager.GetManager().AppendUserConn(userConn)
		fmt.Printf("bindUser success,uid:%d\n", bindUserData.Uid)
	}
}
