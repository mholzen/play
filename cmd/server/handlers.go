package main

import (
	"fmt"
	"io"
	"log"
	"maps"
	"net/http"
	"slices"
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

func ContainerPostHandler(container controls.Container) echo.HandlerFunc {
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
