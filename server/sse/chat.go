package sse

import (
	"wanshow-bingo/utils"
)

var chatHub *Hub

func init() {
	chatHub = NewHub("CHAT")
	go chatHub.Run()
	utils.Debugln("Chat Hub Initialized")
}

func GetChatHub() *Hub {
	return chatHub
}
