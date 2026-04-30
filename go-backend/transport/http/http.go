package transport

import (
	"checker/go-backend/plagiarism/handler"

	"github.com/labstack/echo/v4"
)

type Server struct {
	echo *echo.Echo
	port string
}

func NewServer(h *handler.Handler, port string) *Server {
	e := echo.New()
	e.HideBanner = true

	e.GET("/health", h.Health)

	v1 := e.Group("/api/v1")
	v1.POST("/reference", h.UploadReference)
	v1.POST("/similarity", h.CheckSimilarity)
	v1.POST("/plagiarism", h.CheckPlagiarism)

	return &Server{echo: e, port: port}
}

func (s *Server) Start() error {
	return s.echo.Start(":" + s.port)
}
