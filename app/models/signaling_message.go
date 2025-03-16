package models

type DeviceInfo struct {
}

type IceServerInfo struct {
	Urls []string `json:"urls"`
}

type BaseMessage struct {
	MessageType string `json:"message_type"`
}

type ForwardMessage struct {
	MessageType string      `json:"message_type"`
	ClientId    uint64      `json:"client_id,omitempty"`
	Payload     interface{} `json:"payload"`
}

type RegisterMessage struct {
	MessageType string     `json:"message_type"`
	DevicePort  int        `json:"device_port"`
	DeviceId    string     `json:"device_id"`
	DeviceInfo  DeviceInfo `json:"device_info"`
}

type ConfigMessage struct {
	MessageType string          `json:"message_type"`
	IceServers  []IceServerInfo `json:"ice_servers"`
}

type ConnectMessage struct {
	MessageType string `json:"message_type"`
	DeviceId    string `json:"device_id"`
}

type ClientMessage struct {
	MessageType string      `json:"message_type"`
	ClientId    uint64      `json:"client_id"`
	Payload     interface{} `json:"payload"`
}

type ErrorMessage struct {
	Error string `json:"error"`
}
