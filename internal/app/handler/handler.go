package handler

import (
	"LAB1/internal/app/repository"
	// "net/http"
	// "strconv"
	// "time"

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

func (h *Handler) RegisterHandler(router *gin.Engine) {
	router.GET("/planets", h.GetPlanets)
	router.GET("/planet/:id", h.GetPlanet)
	router.POST("/planet/create-planet", h.CreatePlanet)
	router.DELETE("/planet/:id/delete-planet", h.DeletePlanet)
	router.PUT("/planet/:id/change-planet", h.ChangePlanet)
	router.POST("/planet/:id/add-to-research", h.AddPlanetToResearch)
	router.POST("/planet/:id/create-image", h.UploadImage)


	router.GET("/research/research-cart", h.GetResearchCart)	
	router.GET("/researches", h.GetResearches)
	router.GET("/research/:id", h.GetRsearch)
	router.PUT("/research/:id/change-research", h.ChangeResearch)
	router.PUT("/research/:id/form", h.FormResearch)
	router.PUT("/research/:id/finish", h.ModerateResearch)
	router.DELETE("/research/:id/delete-research", h.DeleteResearch)

	router.DELETE("/planets_research/:planet_id/:research_id", h.DeletePlanetFromResearch)
	router.PUT("/planets_research/:planet_id/:research_id", h.ChangePlanetResearch)

	router.POST("/users/sign-up", h.CreateUser)
	router.GET("/users/profile", h.GetProfile)
	router.PUT("/users/profile", h.ChangeProfile)
	router.POST("/users/sign-in", h.SignIn)
	router.POST("/users/sign-out", h.SignOut)
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.Static("/static", "/home/muka/Рабочий стол/RIP/LAB1/resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": err.Error(),
	})
}

