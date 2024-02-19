package service

import (
	"awesomeProject/manager"
	"awesomeProject/opType"
)

type ChatService struct {
	MainDealer opType.RawMainDeal
}

func (chatService ChatService) Init() {
	chatService.initDealer()
}

func (chatService ChatService) initDealer() {
	chatService.MainDealer.AddOpDealer(BindUserDealer{})
	chatService.MainDealer.AddOpDealer(LogoutUserDealer{})
	chatService.MainDealer.AddOpDealer(SendAllDealer{})
	chatService.MainDealer.AddOpDealer(SendSingleDealer{})
	chatService.MainDealer.AddOpDealer(JoinRoomDealer{})
	chatService.MainDealer.AddOpDealer(CreateRoomDealer{})
	manager.AddRecCallback(chatService)
}

func (chatService ChatService) OnRec(uid int, data []byte) {
	chatService.MainDealer.Deal(uid, data)
}
