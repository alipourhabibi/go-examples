package main

import (
	"fmt"

	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/routers"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/settings"
)

func main() {
	settings.SetUp()

	address := fmt.Sprintf(":%s", settings.ServerSetting.Items.Port)

	r := routers.InitRouters()
	r.Run(address)
}
