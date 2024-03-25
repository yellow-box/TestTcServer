package service

import (
	"awesomeProject/manager"
	"awesomeProject/opType"
)

type ChatService struct {
	MainDealer      opType.RawMainDeal
	BindUserHandler opType.BindUserAction
}

func (chatService ChatService) Init() {
	chatService.initDealer()
}

func (chatService ChatService) initDealer() {
	chatService.MainDealer.AddOpDealer(LogoutUserDealer{})
	chatService.MainDealer.AddOpDealer(SendAllDealer{})
	chatService.MainDealer.AddOpDealer(JoinRoomDealer{})
	chatService.MainDealer.AddOpDealer(HeartBeatDealer{})
	chatService.MainDealer.AddOpDealer(CreateRoomDealer{})
	chatService.BindUserHandler = BindUserCheck{}
	manager.AddRecCallback(chatService)
}

func (chatService ChatService) OnRec(userConn *manager.UserConn, data []byte) {
	chatService.MainDealer.Deal(userConn, chatService.BindUserHandler, data)
}
