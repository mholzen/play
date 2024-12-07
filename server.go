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
	e.GET("/controls", ControlsGetHandler(dialMap))
	dialList := controls.DialList{
		DialMap:     dialMap,
		ChannelList: controls.ChannelList{"r", "g", "b", "a", "w", "uv", "dimmer", "strobe", "speed", "tilt", "pan", "mode"},
	}

	e.GET("/controls/", ControlsGetHandler2(dialList))

	e.GET("/controls/:name", ControlsGetHandler(dialMap))

	postHandler := ControlsPostHandler(dialMap)
	e.POST("/controls/:name", postHandler)
	e.GET("/controls/:name/:value", postHandler)
	e.POST("/controls/:name/:value", postHandler)

	e.GET("/colors/", ColorsGetHandler())
	e.POST("/colors/:name", ColorsPostHandler(dialMap))

	e.GET("/v2/root", ContainerGetHandler(surface))
	e.GET("/v2/root/", ContainerGetHandler(surface))
	e.GET("/v2/root/:name", ContainerGetHandler(surface))
	e.POST("/v2/root/:name/:value", ContainerPostHandler(surface))

	e.Logger.Fatal(e.Start(":1300"))

}
