package models

import "github.com/gorilla/websocket"

type Client struct {
	ClientId uint64
	Conn     *websocket.Conn
	Device   *Device
}
