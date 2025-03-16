package services

import (
	"errors"
	"log"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
	"github.com/simpleyyt/AnCloud-Server/app/models"
)

// ClientService manages all client connections
type ClientService struct {
	// Map of all clients with ClientId as key
	clients     map[uint64]*models.Client
	clientsLock sync.RWMutex
	// Counter for incrementing client IDs
	counter uint64
}

// Init initializes the ClientService
func (s *ClientService) Init() {
	s.clients = make(map[uint64]*models.Client)
}

// CreateClient creates a new client and assigns an ID
func (s *ClientService) CreateClient(conn *websocket.Conn) *models.Client {
	// Atomic increment of counter to ensure thread safety
	clientId := atomic.AddUint64(&s.counter, 1)
	client := &models.Client{
		ClientId: clientId,
		Conn:     conn,
	}

	s.clientsLock.Lock()
	s.clients[clientId] = client
	s.clientsLock.Unlock()

	log.Printf("Client %d connected", clientId)
	return client
}

// RemoveClient removes a client from management
func (s *ClientService) RemoveClient(client *models.Client) error {
	if client == nil {
		return errors.New("client is nil")
	}

	s.clientsLock.Lock()
	defer s.clientsLock.Unlock()

	// If client is associated with a device, remove it from the device
	if client.Device != nil {
		client.Device.RemoveClient(client)
		client.Device = nil
	}

	// Close the websocket connection if it exists
	if client.Conn != nil {
		client.Conn.Close()
	}

	// Remove from the map
	delete(s.clients, client.ClientId)
	log.Printf("Client %d has been removed", client.ClientId)
	return nil
}

// GetClientById gets a client by its ID
func (s *ClientService) GetClientById(clientId uint64) (*models.Client, error) {
	s.clientsLock.RLock()
	defer s.clientsLock.RUnlock()

	client, ok := s.clients[clientId]
	if !ok {
		return nil, errors.New("client not found")
	}
	return client, nil
}

// ConnectToDevice connects a client to a device
func (s *ClientService) ConnectToDevice(client *models.Client, device *models.Device) error {
	if client == nil {
		return errors.New("client is nil")
	}
	if device == nil {
		return errors.New("device is nil")
	}

	// Lock to ensure atomicity of connection operations
	s.clientsLock.Lock()
	defer s.clientsLock.Unlock()

	// Handle existing connection
	if client.Device != nil {
		// 使用Device的RemoveClient方法从原设备中移除客户端
		oldDevice := client.Device
		oldDevice.RemoveClient(client)
	}

	// Establish new connection
	client.Device = device
	device.Clients = append(device.Clients, client)

	log.Printf("Client %d connected to device %s", client.ClientId, device.DeviceId)

	return nil
}

// DisconnectFromDevice disconnects a client from its device
func (s *ClientService) DisconnectFromDevice(client *models.Client) error {
	if client == nil {
		return errors.New("client is nil")
	}

	s.clientsLock.RLock()
	client, ok := s.clients[client.ClientId]
	s.clientsLock.RUnlock()

	if !ok {
		return errors.New("client not found")
	}

	if client.Device == nil {
		return nil // Already disconnected
	}

	// Lock to ensure atomicity of operation
	s.clientsLock.Lock()
	defer s.clientsLock.Unlock()

	device := client.Device

	// 使用Device的RemoveClient方法从设备中移除客户端
	device.RemoveClient(client)

	// 断开连接
	client.Device = nil

	log.Printf("Client %d disconnected from device %s", client.ClientId, device.DeviceId)

	return nil
}

// RemoveAllClients removes all clients from the service
func (s *ClientService) RemoveAllClients() {
	s.clientsLock.Lock()
	defer s.clientsLock.Unlock()

	for clientId, client := range s.clients {
		// If client is associated with a device, remove it from the device
		if client.Device != nil {
			client.Device.RemoveClient(client)
		}

		// Close the websocket connection if it exists
		if client.Conn != nil {
			client.Conn.Close()
		}

		// Remove from the map
		delete(s.clients, clientId)
		log.Printf("Client %d has been removed", clientId)
	}
}

// Singleton implementation
var (
	clientService *ClientService
	clientOnce    sync.Once
)

// GetClientService returns the ClientService singleton
func GetClientService() *ClientService {
	clientOnce.Do(func() {
		clientService = &ClientService{}
		clientService.Init()
	})
	return clientService
}
