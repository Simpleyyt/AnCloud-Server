package controllers

import (
	"encoding/json"
	"errors"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/simpleyyt/AnCloud-Server/app/constants"
	"github.com/simpleyyt/AnCloud-Server/app/models"
	"github.com/simpleyyt/AnCloud-Server/app/services"
)

type SignalingHandler struct {
	deviceHandlers map[string]func(*models.Device, []byte) error
	clientHandlers map[string]func(*models.Client, []byte) error
}

func (h *SignalingHandler) handleRegisterMessage(device *models.Device, data []byte) error {
	log.Println("Register device message:", string(data))
	message := &models.RegisterMessage{}
	err := json.Unmarshal(data, message)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return err
	}
	device.DeviceId = message.DeviceId
	err = services.GetDeviceService().AddDevice(device)
	if err != nil {
		log.Println("Add device error:", err)
		return err
	}
	iceServers := services.GetConfigService().GetIceServers()
	configMessage := models.ConfigMessage{
		MessageType: constants.MESSAGE_TYPE_CONFIG,
		IceServers:  []models.IceServerInfo{iceServers},
	}
	return h.sendMessage(device.Conn, &configMessage)
}

func (h *SignalingHandler) handleConnectMessage(client *models.Client, data []byte) error {
	log.Println("Connect client message:", string(data))
	message := &models.ConnectMessage{}
	err := json.Unmarshal(data, message)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return err
	}

	device, err := services.GetDeviceService().GetDeviceById(message.DeviceId)
	if err != nil {
		log.Println("Get device error:", err)
		return err
	}
	err = services.GetClientService().ConnectToDevice(client, device)
	if err != nil {
		log.Println("Connect to device error:", err)
		return err
	}

	iceServers := services.GetConfigService().GetIceServers()
	configMessage := models.ConfigMessage{
		MessageType: constants.MESSAGE_TYPE_CONFIG,
		IceServers:  []models.IceServerInfo{iceServers},
	}
	return h.sendMessage(client.Conn, &configMessage)
}

func (h *SignalingHandler) handleClientForwardMessage(client *models.Client, data []byte) error {
	log.Println("Forward client message:", string(data))
	message := &models.ForwardMessage{}
	err := json.Unmarshal(data, message)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return err
	}
	clientMessage := models.ClientMessage{
		MessageType: constants.MESSAGE_TYPE_CLIENT_MSG,
		ClientId:    client.ClientId,
		Payload:     message.Payload,
	}
	return h.sendMessage(client.Device.Conn, &clientMessage)
}

func (h *SignalingHandler) handleDeviceForwardMessage(device *models.Device, data []byte) error {
	log.Println("Forward client message:", string(data))
	message := &models.ForwardMessage{}
	err := json.Unmarshal(data, message)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return err
	}
	clientMessage := models.ClientMessage{
		MessageType: constants.MESSAGE_TYPE_DEVICE_MSG,
		Payload:     message.Payload,
	}
	if len(device.Clients) == 0 {
		return h.sendErrorMessage(device.Conn, "No clients connected")
	}
	client := device.GetClientById(message.ClientId)
	if client == nil {
		return h.sendErrorMessage(device.Conn, "Client not found")
	}
	return h.sendMessage(client.Conn, &clientMessage)
}

func (h *SignalingHandler) sendErrorMessage(conn *websocket.Conn, errorString string) error {
	errorMessage := models.ErrorMessage{
		Error: errorString,
	}
	return h.sendMessage(conn, &errorMessage)
}

func (h *SignalingHandler) sendMessage(conn *websocket.Conn, message interface{}) error {
	log.Println("Send message:", message)
	data, err := json.Marshal(message)
	if err != nil {
		log.Println("Marshal error:", err)
		return err
	}
	err = conn.WriteMessage(websocket.TextMessage, data)
	if err != nil {
		return err
	}

	return nil
}

func (h *SignalingHandler) dispatchDeviceMessage(device *models.Device, data []byte) error {
	log.Println("Dispatch device message:", string(data))
	var baseMessage models.RegisterMessage
	err := json.Unmarshal(data, &baseMessage)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return err
	}
	handler, ok := h.deviceHandlers[baseMessage.MessageType]
	if !ok {
		log.Println("Unsupported message type:", baseMessage.MessageType)
		return errors.New("Unsupported message type: " + baseMessage.MessageType)
	}
	return handler(device, data)
}

func (h *SignalingHandler) dispatchClientMessage(client *models.Client, data []byte) error {
	log.Println("Dispatch client message:", string(data))
	var baseMessage models.RegisterMessage
	err := json.Unmarshal(data, &baseMessage)
	if err != nil {
		log.Println("Unmarshal error:", err)
		return err
	}
	handler, ok := h.clientHandlers[baseMessage.MessageType]
	if !ok {
		log.Println("Unsupported message type:", baseMessage.MessageType)
		return errors.New("Unsupported message type: " + baseMessage.MessageType)
	}
	return handler(client, data)
}

func (h *SignalingHandler) HandleClientConnection(conn *websocket.Conn) error {
	var (
		err  error
		data []byte
	)
	defer conn.Close()
	log.Println("Handle client connection")
	client := services.GetClientService().CreateClient(conn)
	for {
		_, data, err = conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Println("Received: ", string(data))
		err = h.dispatchClientMessage(client, data)
		if err != nil {
			h.sendErrorMessage(conn, err.Error())
			break
		}
	}
	services.GetClientService().RemoveClient(client)
	return err
}

func (h *SignalingHandler) HandleDeviceConnection(conn *websocket.Conn) error {
	var (
		err  error
		data []byte
	)
	defer conn.Close()
	log.Println("Handle client connection")
	device := services.GetDeviceService().CreateDevice(conn)
	for {
		_, data, err = conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}

		log.Println("Received: ", string(data))
		err = h.dispatchDeviceMessage(device, data)
		if err != nil {
			h.sendErrorMessage(conn, err.Error())
			break
		}
	}
	services.GetDeviceService().RemoveDevice(device)
	return err
}

func (h *SignalingHandler) Init() {
	h.deviceHandlers = map[string]func(*models.Device, []byte) error{
		constants.MESSAGE_TYPE_REGISTER: h.handleRegisterMessage,
		constants.MESSAGE_TYPE_FORWARD:  h.handleDeviceForwardMessage,
	}
	h.clientHandlers = map[string]func(*models.Client, []byte) error{
		constants.MESSAGE_TYPE_CONNECT: h.handleConnectMessage,
		constants.MESSAGE_TYPE_FORWARD: h.handleClientForwardMessage,
	}
}

var (
	signalingHandler *SignalingHandler
	once             sync.Once
)

func GetSignalingHandler() *SignalingHandler {
	once.Do(func() {
		signalingHandler = &SignalingHandler{}
		signalingHandler.Init()
	})
	return signalingHandler
}
