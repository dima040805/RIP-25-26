package handler

import (
	"fmt"
	"net/http"
	"strconv"

	apitypes "LAB1/internal/app/api_types"
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
	resp := make([]apitypes.PlanetJSON, 0, len(planets))
	for _, r := range planets {
		resp = append(resp, apitypes.PlanetToJSON(r))
	}
	ctx.JSON(http.StatusOK, resp)
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

	ctx.JSON(http.StatusOK, apitypes.PlanetToJSON(*planet))
}

func (h *Handler) CreatePlanet(ctx *gin.Context) {
	var planetJSON apitypes.PlanetJSON
	if err := ctx.BindJSON(&planetJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	planet, err := h.Repository.CreatePlanet(planetJSON)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.Header("Location", fmt.Sprintf("/reactions/%v", planet.ID))
	ctx.JSON(http.StatusCreated, apitypes.PlanetToJSON(planet))
}


func (h *Handler) AddPlanetToResearch(ctx *gin.Context) {
	research, err := h.Repository.GetResearchDraft(h.Repository.GetUserID())
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

	// ctx.Redirect(http.StatusFound, "/planets")
}