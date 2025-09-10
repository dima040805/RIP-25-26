package api

import (
	"LAB1/internal/app/handler"
	"LAB1/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	handler := handler.NewHandler(repo)

	r := gin.Default()
	// добавляем наш html/шаблон
	r.LoadHTMLGlob("/home/muka/Рабочий стол/RIP/LAB1/templates/*")
	r.Static("/static", "/home/muka/Рабочий стол/RIP/LAB1/resources")

	r.GET("/planets", handler.GetPlanets)
	r.GET("/planet/:id", handler.GetPlanet) // вот наш новый обработчик
	r.GET("/research/:id", handler.ResearchHandler)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	log.Println("Server down")
}
