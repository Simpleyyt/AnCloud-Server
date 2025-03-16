package models

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Device struct {
	DeviceId string
	Conn     *websocket.Conn
	Clients  []*Client
	mutex    sync.RWMutex // 添加读写锁保护Clients
}

func (device *Device) AddClient(client *Client) {
	device.mutex.Lock()         // 写操作前加锁
	defer device.mutex.Unlock() // 操作完成后解锁
	device.Clients = append(device.Clients, client)
}

func (device *Device) RemoveClient(client *Client) {
	device.mutex.Lock()         // 写操作前加锁
	defer device.mutex.Unlock() // 操作完成后解锁
	// 使用切片过滤的方式移除客户端
	var newClients []*Client
	for _, c := range device.Clients {
		if c != client {
			newClients = append(newClients, c)
		}
	}
	device.Clients = newClients
}

func (device *Device) GetClientById(clientId uint64) *Client {
	device.mutex.RLock()
	defer device.mutex.RUnlock()
	for _, c := range device.Clients {
		if c.ClientId == clientId {
			return c
		}
	}
	return nil
}
