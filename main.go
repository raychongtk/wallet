package main

import (
	"context"
	"github.com/raychongtk/wallet/util"
	"log"
)

func main() {
	ctx := context.TODO()
	routes, err := injectRoutes(ctx)
	if err != nil {
		log.Fatalln("inject routes failed")
	}
	util.InitializeLogger(false)
	if err := routes.Run(":8080"); err != nil {
		log.Fatalln("start server failed")
	}
}
