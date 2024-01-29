package service

import (
	"awesomeProject/manager"
	"container/list"
)

type ChatService struct {
	registedUsers *list.List
}

func (chatService ChatService) registerUser(user manager.User) {
	chatService.registedUsers.PushBack(user)
}
