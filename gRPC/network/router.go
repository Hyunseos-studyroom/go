package network

import (
	"gRPC/config"
	"gRPC/service"
	"github.com/gin-gonic/gin"
)

type Network struct {
	cfg *config.Config

	service *service.Service

	engine *gin.Engine
}

func NewNetwork(cfg *config.Config, service *service.Service) (*Network, error) {
	r := &Network{cfg: cfg, service: service, engine: gin.New()}

	return r, nil
}

func (s *Network) StartServer() {
	s.engine.Run(":8080")
}
