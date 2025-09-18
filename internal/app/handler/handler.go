package handler

import (
	"LAB1/internal/app/repository"
	// "net/http"
	// "strconv"
	// "time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/planets", h.GetPlanets)
	router.GET("/planet/:id", h.GetPlanet)
	// router.GET("/research/:id", h.ResearchHandler)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.LoadHTMLGlob("/home/muka/Рабочий стол/RIP/LAB1/templates/*")
	router.Static("/static", "/home/muka/Рабочий стол/RIP/LAB1/resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}

// func (h *Handler) GetPlanets(ctx *gin.Context) {
// 	var planets []repository.Planet
// 	var err error

// 	searchQuery := ctx.Query("query")
// 	if searchQuery == "" {
// 		planets, err = h.Repository.GetPlanets()
// 		if err != nil {
// 			logrus.Error(err)
// 		}
// 	} else {
// 		planets, err = h.Repository.GetPlanetsByName(searchQuery)
// 		if err != nil {
// 			logrus.Error(err)
// 		}
// 	}

// 	ctx.HTML(http.StatusOK, "index.html", gin.H{
// 		"planets":       planets,
// 		"query":         searchQuery,
// 		"researchId":    h.Repository.GetResearchId(),
// 		"researchCount": h.Repository.GetResearchCount(1),
// 	})
// }

// func (h *Handler) GetPlanet(ctx *gin.Context) {
// 	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /order/:id)
// 	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
// 	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
// 	if err != nil {
// 		logrus.Error(err)
// 	}

// 	planet, err := h.Repository.GetPlanet(id)
// 	if err != nil {
// 		logrus.Error(err)
// 	}

// 	ctx.HTML(http.StatusOK, "planet.html", gin.H{
// 		"planet": planet,
// 	})
// }

// func (h *Handler) ResearchHandler(ctx *gin.Context) {
// 	idStr := ctx.Param("id")
// 	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
// 	if err != nil {
// 		logrus.Error(err)
// 	}
// 	research := h.Repository.GetResearch(id)
// 	researchPlanets := h.Repository.GetResearchPlanets(id)
// 	ctx.HTML(http.StatusOK, "research.html", gin.H{
// 		"researchPlanets": researchPlanets,
// 		"research":        research,
// 		"count":           h.Repository.GetResearchCount(1),
// 	})
// }
