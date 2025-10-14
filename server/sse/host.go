package sse

import (
	"wanshow-bingo/utils"
)

var hostHub *Hub

func init() {
	hostHub = NewHub("HOST")
	go hostHub.Run()
	utils.Debugln("Host Hub Initialized")
}

func GetHostHub() *Hub {
	return hostHub
}
