package app

import (
	"github.com/gin-gonic/gin"
	"github.com/simpleyyt/AnCloud-Server/app/controllers"
)

func NewRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	apiV1 := r.Group("")
	{
		signalingController := new(controllers.SignalingController)
		apiV1.GET("/register_device", signalingController.RegisterDevice)
		apiV1.GET("/connect_client", signalingController.ConnectClient)
	}

	return r
}
