package network

import (
	"gRPC/config"
	"gRPC/service"
	"github.com/gin-gonic/gin"
)

type Network struct {
	cfg *config.Config

	service *service.Service

	engie *gin.Engine
}

func NewNetwork(cfg *config.Config, service *service.Service) (*Network, error) {
	r := &Network{cfg: cfg, service: service, engie: gin.New()}

	return r, nil
}

func (n *Network) StartServer() error {
	return n.engie.Run(":8080")
}
