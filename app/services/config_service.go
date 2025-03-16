package services

import (
	"github.com/simpleyyt/AnCloud-Server/app/models"
	"github.com/simpleyyt/AnCloud-Server/config"
)

type ConfigService struct{}

var configService = &ConfigService{}

func (s *ConfigService) GetIceServers() models.IceServerInfo {
	return models.IceServerInfo{
		Urls: []string{config.GetConfig().WebRTC.StunServer},
	}
}

func GetConfigService() *ConfigService {
	return configService
}
