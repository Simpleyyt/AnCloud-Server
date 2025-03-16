package services

import (
	"errors"
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/simpleyyt/AnCloud-Server/app/models"
)

type DeviceService struct {
	devices     map[string]*models.Device
	devicesLock sync.RWMutex // Add lock for thread safety
}

// CreateDevice creates a new device with the given websocket connection
func (s *DeviceService) CreateDevice(conn *websocket.Conn) *models.Device {
	device := &models.Device{
		Conn:    conn,
		Clients: make([]*models.Client, 0),
	}

	log.Printf("New device created with connection")
	return device
}

func (s *DeviceService) AddDevice(device *models.Device) error {
	if device.DeviceId == "" {
		return errors.New("device id is empty")
	}

	s.devicesLock.Lock()
	defer s.devicesLock.Unlock()

	s.devices[device.DeviceId] = device
	log.Printf("Device %s registered", device.DeviceId)
	return nil
}

func (s *DeviceService) RemoveDevice(device *models.Device) error {
	if device.DeviceId == "" {
		return errors.New("device id is empty")
	}

	s.devicesLock.Lock()
	defer s.devicesLock.Unlock()

	delete(s.devices, device.DeviceId)
	for _, client := range device.Clients {
		client.Device = nil
	}
	log.Printf("Device %s unregistered", device.DeviceId)
	return nil
}

func (s *DeviceService) GetDeviceById(deviceId string) (*models.Device, error) {
	s.devicesLock.RLock()
	defer s.devicesLock.RUnlock()

	device, ok := s.devices[deviceId]
	if !ok {
		return nil, errors.New("device not found")
	}
	return device, nil
}

func (s *DeviceService) Init() {
	s.devices = make(map[string]*models.Device)
}

var (
	deviceService *DeviceService
	once          sync.Once
)

func GetDeviceService() *DeviceService {
	once.Do(func() {
		deviceService = &DeviceService{}
		deviceService.Init()
	})
	return deviceService
}
