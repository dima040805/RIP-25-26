package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)


func (h *Handler) GetResearchCart(ctx *gin.Context){
	planetsCount := h.Repository.GetResearchCount(h.Repository.GetUserID())

	if planetsCount == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":          "no_draft",
			"reactions_count": planetsCount,
		})
		return
	}

	research, err := h.Repository.CheckCurrentResearchDraft(h.Repository.GetUserID())
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status":          "draft",
		"reactions_count": planetsCount,
		"id":              research.ID,
	})
}

func (h *Handler) ResearchHandler(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}
	researchPlanets, research, err := h.Repository.GetPlanetsResearch(id)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"researchPlanets": researchPlanets,
		"research":        research,
	})
}

func (h *Handler) DeleteResearch(ctx *gin.Context){
	idStr := ctx.Param("id")
	researchId, err := strconv.Atoi(idStr) // так как функция выше возвращает нам строку, нужно ее преобразовать в int
	if err != nil {
		logrus.Error(err)
	}


	err = h.Repository.DeleteCalculation(researchId)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	// ctx.Redirect(http.StatusFound, "/planets")
}


