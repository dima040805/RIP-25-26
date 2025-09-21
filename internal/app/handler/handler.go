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
	router.POST("/planet", h.CreatePlanet)
	router.DELETE("/planet/:id", h.DeletePlanet)
	router.PUT("planet/:id", h.ChangePlanet)
	router.POST("/planet/:id/add-to-research", h.AddPlanetToResearch)


	router.GET("/research/research-cart", h.GetResearchCart)	
	router.GET("/research/:id", h.ResearchHandler)
	router.POST("/research/:id/delete-research", h.DeleteResearch)

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

