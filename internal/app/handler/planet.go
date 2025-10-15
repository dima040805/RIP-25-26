package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	apitypes "LAB1/internal/app/api_types"
	"LAB1/internal/app/ds"
	"LAB1/internal/app/repository"

	"github.com/gin-gonic/gin"
)

// GetPlanets godoc
// @Summary Получить список планет
// @Description Возвращает все планеты или фильтрует по названию
// @Tags planets
// @Produce json
// @Param planet_name query string false "Название планеты для поиска"
// @Success 200 {array} apitypes.PlanetJSON "Список планет"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /planets [get]
func (h *Handler) GetPlanets(ctx *gin.Context) {
	var planets []ds.Planet
	var err error

	searchQuery := ctx.Query("planet_name")
	if searchQuery == "" {
		planets, err = h.Repository.GetPlanets()
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	} else {
		planets, err = h.Repository.GetPlanetsByName(searchQuery)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
	}
	resp := make([]apitypes.PlanetJSON, 0, len(planets))
	for _, r := range planets {
		resp = append(resp, apitypes.PlanetToJSON(r))
	}
	ctx.JSON(http.StatusOK, resp)
}

// GetPlanet godoc
// @Summary Получить планету по ID
// @Description Возвращает информацию о планете по её идентификаторуqqqq
// @Tags planets
// @Produce json
// @Param id path int true "ID планеты"
// @Success 200 {object} apitypes.PlanetJSON "Данные планеты"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 404 {object} map[string]string "Планета не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Router /planet/{id} [get]
func (h *Handler) GetPlanet(ctx *gin.Context) {
	idStr := ctx.Param("id") 
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	planet, err := h.Repository.GetPlanet(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, apitypes.PlanetToJSON(*planet))
}

// CreatePlanet godoc
// @Summary Создать новую планету
// @Description Создает новую планету и возвращает её данные
// @Tags planets
// @Accept json
// @Produce json
// @Param planet body apitypes.PlanetJSON true "Данные новой планеты"
// @Success 201 {object} apitypes.PlanetJSON "Созданная планета"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /planet/create-planet [post]
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

	ctx.Header("Location", fmt.Sprintf("/planets/%v", planet.ID))
	ctx.JSON(http.StatusCreated, apitypes.PlanetToJSON(planet))
}

// DeletePlanet godoc
// @Summary Удалить планету
// @Description Выполняет логическое удаление планеты по ID
// @Tags planets
// @Produce json
// @Param id path int true "ID планеты"
// @Success 200 {object} map[string]string "Статус удаления"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 404 {object} map[string]string "Планета не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /planet/{id}/delete-planet [delete]
func (h *Handler) DeletePlanet(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.DeletePlanet(id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "deleted",
	})
}

// ChangePlanet godoc
// @Summary Изменить данные планеты
// @Description Обновляет информацию о планете по ID
// @Tags planets
// @Accept json
// @Produce json
// @Param id path int true "ID планеты"
// @Param planet body apitypes.PlanetJSON true "Новые данные планеты"
// @Success 200 {object} apitypes.PlanetJSON "Обновленная планета"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 404 {object} map[string]string "Планета не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /planet/{id}/change-planet [put]
func (h *Handler) ChangePlanet(ctx *gin.Context){
	var planetJSON apitypes.PlanetJSON
	if err := ctx.BindJSON(&planetJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	planet, err := h.Repository.ChangePlanet(id, planetJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, apitypes.PlanetToJSON(planet))
}

// AddPlanetToResearch godoc
// @Summary Добавить планету в исследование
// @Description Добавляет планету в черновик исследования пользователя
// @Tags planets
// @Produce json
// @Param id path int true "ID планеты"
// @Success 200 {object} apitypes.ResearchJSON "Исследование с добавленной планетой"
// @Success 201 {object} apitypes.ResearchJSON "Создано новое исследование"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 404 {object} map[string]string "Планета не найдена"
// @Failure 409 {object} map[string]string "Планета уже в исследовании"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /planet/{id}/add-to-research [post]
func (h *Handler) AddPlanetToResearch(ctx *gin.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	research, created, err := h.Repository.GetResearchDraft(userID)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}
	researchId := research.ID

	planetId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	err = h.Repository.AddPlanetToResearch(int(researchId), planetId)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else if errors.Is(err, repository.ErrAlreadyExists) {
			h.errorHandler(ctx, http.StatusConflict, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}
	
	status := http.StatusOK
	
	if created {
		ctx.Header("Location", fmt.Sprintf("/research/%v", research.ID))
		status = http.StatusCreated
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(research)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(status, apitypes.ResearchToJSON(research, creatorLogin, moderatorLogin))
}

// UploadImage godoc
// @Summary Загрузить изображение для планеты
// @Description Загружает изображение для планеты и возвращает обновленные данные
// @Tags planets
// @Accept multipart/form-data
// @Produce json
// @Param id path int true "ID планеты"
// @Param image formData file true "Изображение планеты"
// @Success 200 {object} map[string]interface{} "Статус загрузки и данные планеты"
// @Failure 400 {object} map[string]string "Неверный запрос или файл"
// @Failure 404 {object} map[string]string "Планета не найдена"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /planet/{id}/create-image [post]
func (h *Handler) UploadImage(ctx *gin.Context) {
	planetId, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	planet, err := h.Repository.UploadImage(ctx, planetId, file)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "uploaded",
		"planet": apitypes.PlanetToJSON(planet),
	})
}