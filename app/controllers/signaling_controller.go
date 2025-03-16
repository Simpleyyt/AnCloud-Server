package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域请求
	},
}

func upgradeWebsocket(c *gin.Context) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Websocket upgrade error: %v", err)
	}
	return conn, err
}

type SignalingController struct{}

func (s *SignalingController) RegisterDevice(c *gin.Context) {
	conn, err := upgradeWebsocket(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}
	go GetSignalingHandler().HandleDeviceConnection(conn)
}

func (s *SignalingController) ConnectClient(c *gin.Context) {
	conn, err := upgradeWebsocket(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to upgrade to WebSocket"})
		return
	}
	go GetSignalingHandler().HandleClientConnection(conn)
}
