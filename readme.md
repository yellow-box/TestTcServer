# 简易聊天室服务端

功能：为在聊天室内，并且在线的用户转发文本消息，带有心跳检测功能。
默认监听的端口是12345。

##

UserConn 一个用户套接字的抽象，记录换一个套接字和用户uid
ChatRoom 一个聊天室的抽象 就记录了房间id和房间内的所有uid
userConnManager 所有UserConn的管理类
ChatRoomManager 所有聊天室的管理类
opType 服务端和客户端之间发送消息，每个消息都要带有一个opType，表示该消息的类型,这里的定义了各种类型
ChatService 注册了各种opType的处理函数，当收到客户端消息后，会转发到 ChatService.OnRec,然后根据opType转发到相应的OpTypeDealer

RoomMsg 表示一条聊天室消息
Rsp 给客户端的回包 结构体

## 监听流程

start.go 下的 startListen开始监听当前套接字。

当和客户端的套接字建立链接后，会立即启动 定时心跳包检测，以及通过UserConn.readData()监听该套接字的输入流，
在收到信息后，会在UserConn.onRec接受消息进行处理，然后会转发到ChatService的OnRec，此处会根据注册的OpTypeDealer进行处理。

收到某个用户 opType为 SEND_MESSAGE_ALL 的 消息后，会给除了发送者外房间内所有人发送 OpType为 PUSH_MESSAGE 的消息，为消息发送者发送一个相同OpType的回包（SEND_MESSAGE_ALL）

给特定用户发送消息，是通过uid 在UserConnManager中根据uid拿到userConn，然后再想userConn写入数据

发送消息时，会先发送消息的长度，然后再发送消息的内容，以确保接收方明确接受的数据的长度。