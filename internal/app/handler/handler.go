package handler

import (
	"LAB1/internal/app/repository"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/swaggo/gin-swagger"
	`github.com/swaggo/files`
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
	router.Use(CORSMiddleware())
	
	api := router.Group("/api/v1")

	unauthorized := api.Group("/")
	unauthorized.POST("/users/sign-up", h.SignUp)	
	unauthorized.GET("/planets", h.GetPlanets)
	unauthorized.GET("/planet/:id", h.GetPlanet)
	unauthorized.POST("/users/sign-in", h.SignIn)
	
    unauthorized.PUT("/research/:id/radius", h.UpdatePlanetRadius)


	optionalauthorized := api.Group("/")
	optionalauthorized.Use(h.WithOptionalAuthCheck())
	optionalauthorized.GET("/research/research-cart", h.GetResearchCart)	


	authorized := api.Group("/")
	authorized.Use(h.ModeratorMiddleware(false))

	authorized.POST("/planet/create-planet", h.CreatePlanet)
	authorized.DELETE("/planet/:id/delete-planet", h.DeletePlanet)
	authorized.PUT("/planet/:id/change-planet", h.ChangePlanet)
	authorized.POST("/planet/:id/add-to-research", h.AddPlanetToResearch)
	authorized.POST("/planet/:id/create-image", h.UploadImage)

	authorized.GET("/researches", h.GetResearches)
	authorized.GET("/research/:id", h.GetRsearch)
	authorized.PUT("/research/:id/change-research", h.ChangeResearch)
	authorized.PUT("/research/:id/form", h.FormResearch)
	authorized.DELETE("/research/:id/delete-research", h.DeleteResearch)

	authorized.DELETE("/planets_research/:planet_id/:research_id", h.DeletePlanetFromResearch)
	authorized.PUT("/planets_research/:planet_id/:research_id", h.ChangePlanetResearch)

	authorized.GET("/users/:login/profile", h.GetProfile)
	authorized.PUT("/users/:login/profile", h.ChangeProfile)
	authorized.POST("/users/sign-out", h.SignOut)

	moderator := api.Group("/")
	moderator.Use(h.ModeratorMiddleware(true))
	moderator.PUT("/research/:id/finish", h.ModerateResearch)


	swaggerURL := ginSwagger.URL("/swagger/doc.json")
	router.Any("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, swaggerURL))
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
}

func (h *Handler) RegisterStatic(router *gin.Engine) {
	router.Static("/static", "/home/muka/Рабочий стол/RIP/LAB1/resources")
}

func (h *Handler) errorHandler(ctx *gin.Context, errorStatusCode int, err error) {
	logrus.Error(err.Error())
	
	var errorMessage string
	switch {
	case errors.Is(err, repository.ErrNotFound):
		errorMessage = "Не найден"
	case errors.Is(err, repository.ErrAlreadyExists):
		errorMessage = "Уже существует"
	case errors.Is(err, repository.ErrNotAllowed):
		errorMessage = "Доступ запрещен"
	case errors.Is(err, repository.ErrNoDraft):
		errorMessage = "Черновик не найден"
	default:
		errorMessage = err.Error()
	}
	
	ctx.JSON(errorStatusCode, gin.H{
		"status":      "error",
		"description": errorMessage,
	})
}