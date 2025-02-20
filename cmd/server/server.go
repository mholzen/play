package main

import (
	"path/filepath"

	"github.com/mholzen/play-go/controls"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartServer(surface controls.Container) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Adjust the static file path to be relative to the project root
	publicPath := filepath.Join("..", "..", "public")
	e.Static("/", publicPath)

	v2root := e.Group("/api/v2/root")
	v2root.GET("", ContainerGetHandler(surface))
	// v2root.GET("/", ContainerGetHandler(surface))
	v2root.GET("/*", ContainerGetHandler(surface))
	v2root.POST("/*", ContainerPostHandler(surface))

	e.Logger.Fatal(e.Start(":1300"))
}
