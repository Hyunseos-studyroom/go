package main

import (
	"flag"
	"fmt"
	"gRPC/cmd"
	"gRPC/config"
)

var configFlag = flag.String("config", "config.toml", "config path")

func main() {
	flag.Parse()

	fmt.Println(*configFlag)

	cfg := config.NewConfig(*configFlag)

	cmd.NewApp(cfg)
}
