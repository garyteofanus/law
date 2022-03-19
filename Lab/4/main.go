package main

import (
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	_ "github.com/lib/pq"
	"net/http"
	"strconv"
)

func main() {
	// Init router
	router := echo.New()
	router.HideBanner = true

	// Setup middleware
	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format:           `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
		CustomTimeFormat: "",
	}))
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		MaxAge:       300,
	}))
	router.Use(middleware.Recover())

	// Init database
	db, err := sqlx.Connect("postgres", "user=lab4 password=lab4 dbname=lab4 sslmode=disable TimeZone=Asia/Jakarta")
	if err != nil {
		log.Fatalf("cannot start database: %v", err)
	}

	// Setup module
	repository := NewRepository(db)
	service := NewService(repository)
	handler := NewHandler(service)

	// Setup routes
	router.GET("/", handler.indexTask)
	router.POST("/", handler.createTask)
	taskWithID := router.Group("/:id", func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			pathParam := c.Param("id")
			if pathParam == "" {
				return echo.NewHTTPError(http.StatusBadRequest, "id is required")
			}

			id, err := strconv.Atoi(pathParam)
			if err != nil {
				return c.JSON(http.StatusUnprocessableEntity, response{
					Error: "Cannot convert ID to int",
				})
			}
			c.Set("id", id)

			return next(c)
		}
	})
	{
		taskWithID.GET("", handler.viewTask)
		taskWithID.PUT("", handler.updateTask)
		taskWithID.DELETE("", handler.deleteTask)
	}

	router.POST("/files", handler.uploadFile)

	// Start server
	if err := router.Start(":8080"); err != nil {
		log.Fatalf("cannot start server: %v", err)
	}
}

