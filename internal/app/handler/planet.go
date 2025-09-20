package handler

import (
	"net/http"
	"strconv"

	"LAB1/internal/app/ds"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func (h *Handler) GetPlanets(ctx *gin.Context) {
	var planets []ds.Planet
	var err error
	creatorID := h.Repository.GetUser()

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		planets, err = h.Repository.GetPlanets()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			logrus.Error(err)
			return
		}
	} else {
		planets, err = h.Repository.GetPlanetsByName(searchQuery)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			logrus.Error(err)
			return
		}
	}
	currentResearch, _ := h.Repository.CheckCurrentResearchDraft(creatorID)

	ctx.HTML(http.StatusOK, "planets.html", gin.H{
		"planets":       planets,
		"researchCount": h.Repository.GetResearchCount(),
		"query":         searchQuery,
		"researchId":    currentResearch.ID,
	})
}

func (h *Handler) GetPlanet(ctx *gin.Context) {
	idStr := ctx.Param("id") 
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	planet, err := h.Repository.GetPlanet(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		logrus.Error(err)
		return
	}

	ctx.HTML(http.StatusOK, "planet.html", gin.H{
		"planet": planet,
	})
}


func (h *Handler) AddPlanetToResearch(ctx *gin.Context) {
	research, err := h.Repository.GetResearchDraft(h.Repository.GetUser())
	researchId := research.ID
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	planetId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.AddPlanetToResearch(int(researchId), planetId)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Redirect(http.StatusFound, "/planets")
}