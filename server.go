package main

import (
	"fmt"
	"log"
	"maps"
	"net/http"
	"slices"
	"strconv"

	"github.com/mholzen/play-go/controls"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func ControlsGetHandler(dialMap controls.DialMap) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")
		if name == "" {
			return c.JSON(http.StatusOK, dialMap)
		}
		dial, ok := dialMap[name]
		if !ok {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Cannot find dial name '%s'", name))
		}
		return c.String(http.StatusOK, fmt.Sprintf("%d", dial.Value))
	}
}

func ControlsGetHandler2(dialMap controls.DialList) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")
		if name == "" {
			return c.JSON(http.StatusOK, dialMap)
		}
		dial, ok := dialMap.DialMap[name]
		if !ok {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Cannot find dial name '%s'", name))
		}
		return c.String(http.StatusOK, fmt.Sprintf("%d", dial.Value))
	}
}

func ControlsPostHandler(dialMap controls.DialMap) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")
		dial, ok := dialMap[name]
		if !ok {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Cannot find dial name '%s'", name))
		}

		value := c.Param("value")

		b, err := strconv.Atoi(value)
		if err != nil {
			return err
		}

		dial.SetValue(byte(b))

		return c.String(http.StatusOK, fmt.Sprintf("%d", dial.Value))
	}
}

func ContainerPostHandler(container controls.Container) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")
		item, err := container.GetItem(name)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Error finding emitter named '%s'", name))
		}

		emitter, ok := item.(controls.EmitterI)
		if !ok {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Item is not an emitter (item: '%s')", name))
		}

		value := c.Param("value")
		emitter.SetValue(value)

		return c.String(http.StatusOK, fmt.Sprintf("%v", emitter.GetValue()))
	}
}

func ContainerGetHandler(container controls.Container) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")
		if name == "" {
			path := c.Path()
			if path[len(path)-1] == '/' {
				return c.JSON(http.StatusOK, slices.Sorted(maps.Keys(container.Items())))
			} else {
				return c.JSON(http.StatusOK, container.Items())
			}
		}
		// request path ends with /
		item, err := container.GetItem(name)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Error finding emitter named '%s'", name))
		}

		return c.JSON(http.StatusOK, item)
	}
}

func ColorsPostHandler(dialMap controls.DialMap) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")
		color, ok := controls.AllColors[name]
		if !ok {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Cannot find color name '%s'", name))
		}

		dialMap.SetValues(color.Values())

		return c.String(http.StatusOK, fmt.Sprintf("%v", color))
	}
}

func ColorsGetHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, controls.AllColors)
	}
}

func StartServer(surface controls.Container) {
	controls.LoadColors()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "public")

	item0, err := surface.GetItem("0")
	if err != nil {
		log.Fatalf("Error channel map: %v", err)
	}
	dialMap := item0.(controls.DialMap)
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

	e.Logger.Fatal(e.Start(":1300"))

}
