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

func (h *Handler) AddPlanetToResearch(ctx *gin.Context) {
	research, created, err := h.Repository.GetResearchDraft(h.Repository.GetUserID())
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