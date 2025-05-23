package main

import (
	"log"

	"github.com/mholzen/play-go/controls"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func StartServer(surface controls.Container) {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "public")

	item0, err := surface.GetItem("0")
	if err != nil {
		log.Fatalf("Error channel map: %v", err)
	}
	dialMap := item0.(*controls.ObservableDialMap)

	dialList := controls.DialList{
		DialMap:     dialMap,
		ChannelList: controls.ChannelList{"r", "g", "b", "a", "w", "uv", "dimmer", "strobe", "speed", "tilt", "pan", "mode"},
	}

	// Create a group with prefix "/controls"
	controls := e.Group("/api/controls")

	getHandler := ControlsGetHandler(dialList)
	controls.GET("/", getHandler)
	controls.GET("/:name", getHandler)

	postHandler := ControlsPostHandler(dialMap)
	controls.POST("/:name", postHandler)
	controls.GET("/:name/:value", postHandler)
	controls.POST("/:name/:value", postHandler)

	// Create a group with prefix "/colors"
	colors := e.Group("/api/colors")
	colors.GET("/", ColorsGetHandler())
	colors.POST("/:name", ColorsPostHandler(dialMap))

	// Create a group with prefix "/v2/root"
	v2root := e.Group("/api/v2/root")
	v2root.GET("", ContainerGetHandler(surface))
	v2root.GET("/", ContainerGetHandler(surface))
	v2root.GET("/*", ContainerGetHandler(surface))
	v2root.POST("/*", ContainerPostHandler(surface))

	e.Logger.Fatal(e.Start(":1300"))
}
