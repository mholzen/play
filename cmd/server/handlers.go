package main

import (
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mholzen/play-go/controls"
)

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

// PathResolve resolves a path within a container by traversing the path segments
func PathResolve(container controls.Container, path string) (controls.Item, string, error) {
	if path == "" {
		return container, path, nil
	}

	segments := strings.Split(path, "/")

	current := container
	for i, segment := range segments {
		lead := strings.Join(segments[:i], "/")
		remainder := strings.Join(segments[i+1:], "/")
		if segment == "" {
			continue
		}

		child, err := current.GetItem(segment)
		if err != nil {
			return nil, remainder, fmt.Errorf("failed to resolve path segment '%s': %w", segment, err)
		}
		if i == len(segments)-1 {
			return child, "", nil
		}
		if newContainer, ok := child.(controls.Container); !ok {
			return nil, remainder, fmt.Errorf("path segment '%s' is not a container, got %T", lead, child)
		} else {
			current = newContainer
		}
	}

	return current, "", nil
}

func ContainerPostHandler2(container controls.Container) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Param("*")
		if path == "" {
			return echo.NewHTTPError(http.StatusBadRequest, "path is required")
		}

		item, _, err := PathResolve(container, path)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Error finding control named '%s': %v", path, err))
		}

		body := c.Request().Body
		defer body.Close()
		bytes, err := io.ReadAll(body)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("request body is required as the value: %v", err))
		}
		value := string(bytes)

		switch control := item.(type) {
		case controls.Control:
			control.SetValueString(value)
			log.Printf("Control '%s' updated to '%s'", path, value)
		case controls.Container:
			return echo.NewHTTPError(http.StatusMethodNotAllowed, fmt.Sprintf("Cannot post directly to amake test container '%s'", path))
		default:
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Item '%s' is not a control (got '%T')", path, item))
		}

		return c.JSON(http.StatusOK, item)
	}
}
