package main

import (
	"fmt"
	"net/http"
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

func StartServer(surface controls.DialMap) {
	controls.LoadColors()

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/", "public")

	e.GET("/controls", ControlsGetHandler(surface))
	dialList := controls.DialList{
		DialMap:     surface,
		ChannelList: controls.ChannelList{"r", "g", "b", "a", "w", "uv", "dimmer", "strobe", "speed", "tilt", "pan", "mode"},
	}

	e.GET("/controls/", ControlsGetHandler2(dialList))

	e.POST("/controls/:name", ControlsPostHandler(surface))
	e.POST("/controls/:name/:value", ControlsPostHandler(surface))
	e.GET("/controls/:name/:value", ControlsPostHandler(surface))

	e.POST("/colors/:name", ColorsPostHandler(surface))
	e.GET("/colors/:name", ColorsPostHandler(surface))
	e.GET("/colors/", ColorsGetHandler())

	e.GET("/controls/:name", ControlsGetHandler(surface))

	e.Logger.Fatal(e.Start(":1300"))
}
