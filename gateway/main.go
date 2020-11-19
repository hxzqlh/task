package main

import (
	"log"
	"task/gateway/plugins/auth"

	"github.com/micro/micro/v2/cmd"
	"github.com/micro/micro/v2/plugin"
)

func main() {
	// 注册auth插件
	err := plugin.Register(auth.NewPlugin())
	if err != nil {
		log.Fatal("auth register")
	}

	cmd.Init()
}
