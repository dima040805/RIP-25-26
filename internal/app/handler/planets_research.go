package handler


import (
	apitypes "LAB1/internal/app/api_types"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

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
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	creatorLogin, moderatorLogin, err := h.Repository.GetModeratorAndCreatorLogin(research)
	if err != nil {
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, apitypes.ResearchToJSON(research, creatorLogin, moderatorLogin))
}

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
		h.errorHandler(ctx, http.StatusInternalServerError, err)
		return
	}

	ctx.JSON(http.StatusOK, apitypes.PlanetsResearchToJSON(planetResearch))

}
