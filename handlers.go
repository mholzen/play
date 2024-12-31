package main

import (
	"fmt"
	"log"
	"maps"
	"net/http"
	"slices"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/mholzen/play-go/controls"
)

func ControlsGetHandler(dialList controls.DialList) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")
		if name == "" {
			return c.JSON(http.StatusOK, dialList)
		}
		dial, ok := (*dialList.DialMap.Dials)[name]
		if !ok {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Cannot find dial name '%s'", name))
		}
		return c.String(http.StatusOK, fmt.Sprintf("%d", dial.Value))
	}
}

func ControlsPostHandler(dialMap *controls.ObservableDialMap) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")
		dial, ok := (*dialMap.Dials)[name]
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

func ContainerGetHandler(container controls.Container) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")
		path := c.Path()
		endsWithSlash := path[len(path)-1] == '/'

		// log.Printf("ContainerGetHandler: name: %s, path: %s, endsWithSlash: %t", name, path, endsWithSlash)
		if name == "" {
			if endsWithSlash {
				return c.JSON(http.StatusOK, slices.Sorted(maps.Keys(container.Items())))
			} else {
				return c.JSON(http.StatusOK, container.Items())
			}
		}

		item, err := container.GetItem(name)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Error finding control named '%s'", name))
		}

		if endsWithSlash {
			// Try to convert item to a map/array-like structure
			if container, ok := item.(controls.Container); ok {
				return c.JSON(http.StatusOK, slices.Sorted(maps.Keys(container.Items())))
			}
		}

		return c.JSON(http.StatusOK, item)
	}
}

func ContainerPostHandler(container controls.Container) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")
		item, err := container.GetItem(name)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Error finding control named '%s'", name))
		}

		// convert this to a control -- make observable dialmap a subcase
		switch control := item.(type) {
		case *controls.ObservableDialMap:
			channel := c.Param("channel")
			value := c.Param("value")
			b, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			control.SetChannelValue(channel, byte(b))
			log.Printf("Dial Map updated: %s:%s", name, value)

		case controls.Container:
			channel := c.Param("channel")
			value := c.Param("value")
			b, err := strconv.Atoi(value)
			if err != nil {
				return err
			}
			item, err := control.GetItem(channel)
			if err != nil {
				return err
			}
			switch item := item.(type) {
			case *controls.FloatDial:
				item.SetValue(float64(b))
				log.Printf("Float Dial updated: %s:%s", name, value)
			case *controls.NumericDial:
				item.SetValue(byte(b))
				log.Printf("Numeric Dial updated: %s:%s", name, value)
			default:
				return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Container %s item '%s' is not settable ('%v')", name, channel, item))
			}

			log.Printf("Dial Map updated: %s:%s", name, value)

		case controls.Control:
			value := c.Param("value")
			control.SetValueString(value)
			log.Printf("Control updated: %s", value)

		default:
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Item '%s' is not a dialmap, container or control ('%v')", name, item))
		}

		return c.JSON(http.StatusOK, item)
	}
}

func ColorsPostHandler(dialMap *controls.ObservableDialMap) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := c.Param("name")
		color, ok := controls.AllColors[name]
		if !ok {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Cannot find color name '%s'", name))
		}

		dialMap.SetValue(color.Values())

		return c.String(http.StatusOK, fmt.Sprintf("%v", color))
	}
}

func ColorsGetHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, controls.AllColors)
	}
}
