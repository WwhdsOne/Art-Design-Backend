package bootstrap

import (
	"Art-Design-Backend/pkg/ws"
)

func InitWebSocketHub() *ws.Hub {
	hub := ws.NewHub()
	go hub.Run()
	return hub
}
