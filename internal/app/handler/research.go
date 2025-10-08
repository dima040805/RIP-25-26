package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"LAB1/internal/app/api_types"
	"LAB1/internal/app/ds"
	"LAB1/internal/app/repository"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// GetResearches godoc
// @Summary Получить список исследований
// @Description Возвращает исследования с возможностью фильтрации по датам и статусу
// @Tags researches
// @Produce json
// @Param from-date query string false "Начальная дата (YYYY-MM-DD)"
// @Param to-date query string false "Конечная дата (YYYY-MM-DD)"
// @Param status query string false "Статус исследования"
// @Success 200 {array} apitypes.ResearchJSON "Список исследований"
// @Failure 400 {object} map[string]string "Неверный формат даты"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /researches [get]
func (h *Handler) GetResearches(ctx *gin.Context) {
	fromDate := ctx.Query("from-date")
	var from = time.Time{}
	var to = time.Time{}
	if fromDate != "" {
		from1, err := time.Parse("2006-01-02", fromDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		from = from1
	}

	toDate := ctx.Query("to-date")
	if toDate != "" {
		to1, err := time.Parse("2006-01-02", toDate)
		if err != nil {
			h.errorHandler(ctx, http.StatusBadRequest, err)
			return
		}
		to = to1
	}

	status := ctx.Query("status")

	researches, err := h.Repository.GetResearches(from, to, status)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	researches = h.filterResearchesByAuth(researches, ctx)




	resp := make([]apitypes.ResearchJSON, 0, len(researches))
	for _, c := range researches {
		creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(c)
		if err != nil {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
			return
		}
		resp = append(resp, apitypes.ResearchToJSON(c, creatorLogin, moderatorLogin))
	}
	ctx.JSON(http.StatusOK, resp)
}

// GetResearchCart godoc
// @Summary Получить корзину исследования
// @Description Возвращает информацию о текущем черновике исследования пользователя
// @Tags researches
// @Produce json
// @Success 200 {object} map[string]interface{} "Данные корзины исследования"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /research/research-cart [get]
func (h *Handler) GetResearchCart(ctx *gin.Context){
	userID, err := getUserID(ctx)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}
	planetsCount := h.Repository.GetResearchCount(userID)

	if planetsCount == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"status":          "no_draft",
			"planets_count": planetsCount,
		})
		return
	}

	research, err := h.Repository.CheckCurrentResearchDraft(userID)
	if err != nil {
		if errors.Is(err, repository.ErrNotAllowed) {
			h.errorHandler(ctx, http.StatusUnauthorized, err)
		} else if errors.Is(err, repository.ErrNoDraft) {
			ctx.JSON(http.StatusOK, gin.H{
				"status":          "no_draft",
				"planets_count": 0,
			})
		} else {
			h.errorHandler(ctx, http.StatusInternalServerError, err)
		}
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"id":          research.ID,
		"planets_count": h.Repository.GetResearchCount(research.CreatorID),
	})
}

// GetRsearch godoc
// @Summary Получить исследование по ID
// @Description Возвращает полную информацию об исследовании включая планеты
// @Tags researches
// @Produce json
// @Param id path int true "ID исследования"
// @Success 200 {object} map[string]interface{} "Данные исследования с планетами"
// @Failure 400 {object} map[string]string "Неверный ID"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Исследование не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /research/{id} [get]
func (h *Handler) GetRsearch(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	planets, research, err := h.Repository.GetResearchPlanets(id)
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

	resp := make([]apitypes.PlanetJSON, 0, len(planets))
	for _, r := range planets {
		resp = append(resp, apitypes.PlanetToJSON(r))
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(research)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	planetsResearch, _ := h.Repository.GetPlanetsResearches(research.ID)
	
	resp2 := make([]apitypes.PlanetsResearchJSON, 0, len(planetsResearch))
	for _, r := range planetsResearch {
		resp2 = append(resp2, apitypes.PlanetsResearchToJSON(r))
	}

	ctx.JSON(http.StatusOK, gin.H{
		"research": apitypes.ResearchToJSON(research, creatorLogin, moderatorLogin),
		"planets":   resp,
		"planetsResearch": resp2,
	})
}

// FormResearch godoc
// @Summary Сформировать исследование
// @Description Переводит исследование в статус "formed"
// @Tags researches
// @Produce json
// @Param id path int true "ID исследования"
// @Success 200 {object} apitypes.ResearchJSON "Сформированное исследование"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Исследование не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /research/{id}/form [put]
func (h *Handler) FormResearch(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	status := "formed"

	research, err := h.Repository.FormResearch(id, status)
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

// ChangeResearch godoc
// @Summary Изменить исследование
// @Description Обновляет данные исследования
// @Tags researches
// @Accept json
// @Produce json
// @Param id path int true "ID исследования"
// @Param research body apitypes.ResearchJSON true "Новые данные исследования"
// @Success 200 {object} apitypes.ResearchJSON "Обновленное исследование"
// @Failure 400 {object} map[string]string "Неверные данные"
// @Failure 404 {object} map[string]string "Исследование не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /research/{id}/change-research [put]
func (h *Handler) ChangeResearch(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	var researchJSON apitypes.ResearchJSON
	if err := ctx.BindJSON(&researchJSON); err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	research, err := h.Repository.ChangeResearch(id, researchJSON)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			h.errorHandler(ctx, http.StatusNotFound, err)
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

// DeleteResearch godoc
// @Summary Удалить исследование
// @Description Выполняет логическое удаление исследования
// @Tags researches
// @Produce json
// @Param id path int true "ID исследования"
// @Success 200 {object} map[string]string "Статус удаления"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Исследование не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /research/{id}/delete-research [delete]
func (h *Handler) DeleteResearch(ctx *gin.Context){
	idStr := ctx.Param("id")
	researchId, err := strconv.Atoi(idStr)
	if err != nil {
		h.errorHandler(ctx, http.StatusBadRequest, err)
		return
	}

	status := "deleted"
	
	_, err = h.Repository.FormResearch(researchId, status)
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

	ctx.JSON(http.StatusOK, gin.H{"message": "Research deleted"})
}

// ModerateResearch godoc
// @Summary Модерировать исследование
// @Description Изменяет статус исследования (только для модераторов)
// @Tags researches
// @Accept json
// @Produce json
// @Param id path int true "ID исследования"
// @Param status body apitypes.StatusJSON true "Новый статус"
// @Success 200 {object} apitypes.ResearchJSON "Результат модерации"
// @Failure 400 {object} map[string]string "Неверный запрос"
// @Failure 403 {object} map[string]string "Доступ запрещен"
// @Failure 404 {object} map[string]string "Исследование не найдено"
// @Failure 500 {object} map[string]string "Внутренняя ошибка сервера"
// @Security ApiKeyAuth
// @Router /research/{id}/finish [put]
func (h *Handler) ModerateResearch(ctx *gin.Context) {
    userID, err := getUserID(ctx)
    if err != nil {
        h.errorHandler(ctx, http.StatusBadRequest, err)
        return
    }

    idStr := ctx.Param("id")
    id, err := strconv.Atoi(idStr)
    if err != nil {
        h.errorHandler(ctx, http.StatusBadRequest, err)
        return
    }

    var statusJSON apitypes.StatusJSON
    if err := ctx.BindJSON(&statusJSON); err != nil {
        h.errorHandler(ctx, http.StatusBadRequest, err)
        return
    }

    user, err := h.Repository.GetUserByID(userID)
    if err != nil {
        if errors.Is(err, repository.ErrNotFound) {
            h.errorHandler(ctx, http.StatusNotFound, err)
        } else {
            h.errorHandler(ctx, http.StatusInternalServerError, err)
        }
        return
    }
    
    if !user.IsModerator {
        h.errorHandler(ctx, http.StatusForbidden, errors.New("требуются права модератора"))
        return
    }

    research, err := h.Repository.ModerateResearch(id, statusJSON.Status, userID)
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

func (h *Handler) filterResearchesByAuth(researches []ds.Research, ctx *gin.Context) []ds.Research {
	userID, err := getUserID(ctx)
	if err != nil {
		return []ds.Research{}
	}

	user, err := h.Repository.GetUserByID(userID)
	if err == repository.ErrNotFound {
		return []ds.Research{}
	}
	if err != nil {
		return []ds.Research{}
	}

	if user.IsModerator {
		return researches
	}

	var userResearches []ds.Research
    for _, research := range researches {
        fmt.Println(research.ID)
        if research.CreatorID == userID {
            userResearches = append(userResearches, research)
        }
    }
    
    return userResearches

}

func (h *Handler) hasAccessToResearch(creatorID uuid.UUID, ctx *gin.Context) bool {
	userID, err := getUserID(ctx)
	if err != nil {
		return false
	}

	user, err := h.Repository.GetUserByID(userID)
	if err == repository.ErrNotFound {
		return false
	}
	if err != nil {
		return false
	}

	return creatorID == userID || user.IsModerator
}