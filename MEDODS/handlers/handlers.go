package handlers

import (
	"MEDODS/pkg/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	service *service.Service
}

func New(service *service.Service) *Handlers {
	return &Handlers{service: service}
}

// Initialize gin and EndPoints
func (h *Handlers) InitRoutes() *gin.Engine {
	router := gin.New()

	//Load All HTML files
	router.LoadHTMLGlob("../templates/*")
	fs := http.FileSystem(http.Dir("../static"))
	router.StaticFS("static/", fs)

	//All handlers
	home := router.Group("/", h.identifyUser)
	{
		home.GET("/", h.homeGet)
	}
	auth := router.Group("/auth", h.identifyUser)
	{
		//Registration
		auth.GET("/sign-up", h.signUpGet)
		auth.POST("/sign-up", h.signUpPost)

		//Authorization
		auth.GET("/sign-in", h.signInGet)
		auth.POST("/sign-in", h.signInPost)

	}
	refresh := router.Group("/refresh-tokens")
	{
		refresh.POST("/", h.refreshTokens)
	}

	return router
}
