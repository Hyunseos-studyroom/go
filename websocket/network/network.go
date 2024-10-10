package network

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
)

type Network struct {
	engin *gin.Engine
}

func NewServer() *Network {
	n := &Network{engin: gin.New()}

	n.engin.Use(gin.Logger())
	n.engin.Use(gin.Recovery())
	n.engin.Use(cors.New(cors.Config{
		AllowWebSockets:  true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	r := NewRoom()
	go r.RunInit()

	n.engin.GET("/room", r.SocketServe)
	return n
}

func (n *Network) StartServer() error {
	log.Printf("Starting server")
	return n.engin.Run(":8080")
}
