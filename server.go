package main

import (
	"github.com/mholzen/play-go/controls"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartServer(surface controls.Container) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "public")

	v2root := e.Group("/api/v2/root")
	v2root.GET("", ContainerGetHandler(surface))
	v2root.GET("/", ContainerGetHandler(surface))
	v2root.GET("/:name", ContainerGetHandler(surface))
	v2root.GET("/:name/", ContainerGetHandler(surface))
	v2root.POST("/*", ContainerPostHandler2(surface))

	e.Logger.Fatal(e.Start(":1300"))
}
