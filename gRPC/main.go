package main

import (
	"flag"
	"gRPC/cmd"
	"gRPC/config"
)

var configFlag = flag.String("config", "./config.toml", "config path")

func main() {
	flag.Parse()

	cfg := config.NewConfig(*configFlag)
	cmd.NewApp(cfg)
}
