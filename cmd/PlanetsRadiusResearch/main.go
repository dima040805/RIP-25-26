// RegisterHandler godoc
// @title Astronomy Research API
// @version 1.0
// @description API для управления астрономическими исследованиями планет
// @host localhost:8080
// @BasePath /api/v1
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
package main

import (
	_ "github.com/dima040805/RIP-25-26/docs"
	"github.com/dima040805/RIP-25-26/internal/app/config"
	"github.com/dima040805/RIP-25-26/internal/app/dsn"
	"github.com/dima040805/RIP-25-26/internal/app/handler"
	"github.com/dima040805/RIP-25-26/internal/app/repository"
	"github.com/dima040805/RIP-25-26/internal/pkg"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()

	rep, errRep := repository.NewRepository(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	hand := handler.NewHandler(rep)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}
