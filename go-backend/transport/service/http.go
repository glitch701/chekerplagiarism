package service

import (
	"checker/go-backend/plagiarism/handler"
	"checker/go-backend/plagiarism/service"

	"github.com/labstack/echo/v4"
)

func NewService(router *echo.Echo, svc service.PlagiarismService) {
	routerGroup := router.Group("/api/v1")
	{
		handler.NewHandler(routerGroup, svc)
	}
}
