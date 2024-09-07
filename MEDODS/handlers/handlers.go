package handlers

import (
	"MEDODS/pkg/service"

	"github.com/gin-gonic/gin"
)

type Handlers struct {
	service *service.Service
}

func New(service *service.Service) *Handlers {
	return &Handlers{service: service}
}

// Инициализируем пути по которым будут идти запросы
func (h *Handlers) InitRoutes() *gin.Engine {
	router := gin.New()

	//Запрос на получение Access и Refresh Token
	router.POST("/get_token", h.getToken)
	//Запрос на обновление Refresh Token
	router.POST("/refresh_token", h.refreshToken)

	return router
}
