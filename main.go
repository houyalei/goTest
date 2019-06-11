package main

import (
	"fmt"
	"log"
	"syscall"
	"work/pkg/setting"
	"work/routers"

	"github.com/fvbock/endless"
)

func main() {
	endPoint := fmt.Sprintf(":%d", setting.HTTPPort)

	server := endless.NewServer(endPoint, routers.InitRouter())
	server.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}

}
