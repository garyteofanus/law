package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gitlab.com/garyteofanus/law-assignment/auth"
	"log"
	"net/http"
)

type server struct {
	router *echo.Echo
	port   string

	authHandler auth.Handler
}

func NewServer(router *echo.Echo, port string, authHandler auth.Handler) *server {
	s := &server{
		router: router,
		port:   port,

		authHandler: authHandler,
	}

	// Middlewares
	s.router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{echo.POST},
		AllowHeaders:     []string{echo.HeaderAuthorization, echo.HeaderContentType},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	s.router.Use(middleware.Logger())

	// Routes
	oauthRoute := s.router.Group("/oauth")
	{
		oauthRoute.POST("/token", authHandler.HandleAuth)
		// Protected resource, need to go through auth middleware
		oauthRoute.POST("/resource", authHandler.HandleResource)
	}

	return s
}

func (s *server) Run() {
	if err := http.ListenAndServe(":"+s.port, s.router); err != nil {
		log.Fatalf("port is being used: %v", err)
	}
}
