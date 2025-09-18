package handler

import (
	"LAB1/internal/app/repository"
	"net/http"
	"strconv"
	"time"

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

func (h *Handler) GetPlanets(ctx *gin.Context) {
	var planets []repository.Planet
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		planets, err = h.Repository.GetPlanets()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		planets, err = h.Repository.GetPlanetsByName(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"time":          time.Now().Format("15:04:05"),
		"planets":       planets,
		"query":         searchQuery,
		"researchId":    h.Repository.GetResearchId(),
		"researchCount": h.Repository.GetResearchCount(1),
	})
}

func (h *Handler) GetPlanet(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /order/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}

	planet, err := h.Repository.GetPlanet(id)
	if err != nil {
		logrus.Error(err)
	}

	ctx.HTML(http.StatusOK, "planet.html", gin.H{
		"planet": planet,
	})
}

func (h *Handler) ResearchHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}
	research := h.Repository.GetResearch(id)
	researchPlanets := h.Repository.GetResearchPlanets(id)
	ctx.HTML(http.StatusOK, "research.html", gin.H{
		"researchPlanets": researchPlanets, 
		"research": research,
		"count":           h.Repository.GetResearchCount(1),
	})
}
