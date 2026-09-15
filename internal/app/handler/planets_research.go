package handler

import (
	apitypes "LAB1/internal/app/api_types"
	"LAB1/internal/app/repository"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// DeletePlanetFromResearch godoc
// @Summary Удалить планету из исследования
// @Description Удаляет связь планеты и исследования
// @Tags planets-research
// @Produce json
// @Param planet_id path int true "ID планеты"
// @Param research_id path int true "ID исследования"
// @Success 200 {object} apitypes.ResearchJSON "Обновленное исследование"
// @Failure 400 {object} map[string]string "Неверные ID"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /planets_research/{planet_id}/{research_id} [delete]
func (h *Handler) DeletePlanetFromResearch(ctx *gin.Context) {
	researchId, err := strconv.Atoi(ctx.Param("research_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	planetId, err := strconv.Atoi(ctx.Param("planet_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	research, err := h.Repository.DeletePlanetFromResearch(researchId, planetId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusForbidden, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(research)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, apitypes.ResearchToJSON(research, creatorLogin, moderatorLogin))
}

// ChangePlanetResearch godoc
// @Summary Изменить данные планеты в исследовании
// @Description Обновляет параметры планеты в конкретном исследовании
// @Tags planets-research
// @Accept json
// @Produce json
// @Param planet_id path int true "ID планеты"
// @Param research_id path int true "ID исследования"
// @Param data body apitypes.PlanetsResearchJSON true "Новые данные"
// @Success 200 {object} apitypes.PlanetsResearchJSON "Обновленные данные"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 404 {object} map[string]string "Не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /planets_research/{planet_id}/{research_id} [put]
func (h *Handler) ChangePlanetResearch(ctx *gin.Context) {
	researchId, err := strconv.Atoi(ctx.Param("research_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	planetId, err := strconv.Atoi(ctx.Param("planet_id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var planetResearchJSON apitypes.PlanetsResearchJSON
	if err := ctx.BindJSON(&planetResearchJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	planetResearch, err := h.Repository.ChangePlanetResearch(researchId, planetId, planetResearchJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, apitypes.PlanetsResearchToJSON(planetResearch))
}