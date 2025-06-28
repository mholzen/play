package main

import (
	"fmt"
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
		item, path, err := resolvePathToItem(c, container)
		if err != nil {
			return err
		}

		// endsWithSlash := path[len(path)-1] == '/'
		endsWithSlash := strings.HasSuffix(path, "/")

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

func resolvePathToItem(c echo.Context, container controls.Container) (controls.Item, string, error) {
	path := c.Param("*")
	if path == "" {
		return container, path, nil // the root container is returned
	}

	item, _, err := PathResolve(container, path)
	if err != nil {
		return nil, path, echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Error finding control named '%s': %v", path, err))
	}

	return item, path, nil
}

func ContainerPostHandler(container controls.Container) echo.HandlerFunc {
	return func(c echo.Context) error {
		item, path, err := resolvePathToItem(c, container)
		if err != nil {
			return err
		}

		control, ok := item.(controls.Control)
		if !ok {
			return echo.NewHTTPError(http.StatusMethodNotAllowed, fmt.Sprintf("Item '%s' is not a Control (got '%T')", path, item))
		}

		var data interface{}
		if err := c.Bind(&data); err != nil {
			// this fails is content type is undefined.  it should default to application/json
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("request body is required as the value: %v", err))
		}

		switch value := data.(type) {
		case string:
			control.SetValueString(value)
			log.Printf("control '%s' updated to '%s'", path, value)

		case int:
			stringValue := fmt.Sprintf("%d", value)
			control.SetValueString(stringValue)
			log.Printf("control '%s' updated to '%d'", path, value)

		case float64:
			stringValue := fmt.Sprintf("%d", int(value))
			control.SetValueString(stringValue)
			log.Printf("control '%s' updated to '%f'", path, value)

		case bool:
			stringValue := fmt.Sprintf("%t", value)
			control.SetValueString(stringValue)
			log.Printf("control '%s' updated to '%t'", path, value)
		default:
			return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("request body of type '%T' cannot be converted to a string", value))
		}

		return c.JSON(http.StatusOK, item)
	}
}
