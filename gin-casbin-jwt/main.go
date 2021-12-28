package main

import (
	"fmt"

	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/repo"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/routers"
	"github.com/alipourhabibi/go-examples/gin-casbin-jwt/settings"
)

func main() {
	settings.SetUp()

	repo.AutoMigrateAll()

	address := fmt.Sprintf(":%s", settings.ServerSetting.Items.Port)

	r := routers.InitRouters()
	r.Run(address)
}
