package main

import "websocket/network"

func main() {
	n := network.NewServer()
	n.StartServer()
}
