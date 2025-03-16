package main

import (
	"fmt"

	"github.com/simpleyyt/AnCloud-Server/app"
	"github.com/simpleyyt/AnCloud-Server/config"
)

func main() {
	r := app.NewRouter()
	cfg := config.GetConfig()
	addr := fmt.Sprintf("%s:%d", cfg.Http.Host, cfg.Http.Port)
	r.Run(addr)
}
