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

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"planets": planets,
		"researchCount": h.Repository.GetResearchCount(),
		"query":   searchQuery,
	})
}

func (h *Handler) GetPlanet(ctx *gin.Context) {
	idStr := ctx.Param("id") // получаем id заказа из урла (то есть из /order/:id)
	// через двоеточие мы указываем параметры, которые потом сможем считать через функцию выше
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
