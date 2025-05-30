package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mholzen/play-go/controls"
)

// sanitizePath collapses multiple consecutive slashes into a single slash
func sanitizePath(path string) string {
	re := regexp.MustCompile(`/{2,}`)
	return re.ReplaceAllString(path, "/")
}

func ContainerGetHandler(container controls.Container) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Param("*")
		if path == "" {
			return c.JSON(http.StatusOK, container.Items())
		}

		// Sanitize path by collapsing duplicate slashes
		path = sanitizePath(path)
		segments := strings.Split(path, "/")
		item, err := controls.ContainerFollowPath(container, segments)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Errorf("error following path '%s': %w", path, err))
		}

		return c.JSON(http.StatusOK, item)
	}
}

func ContainerPostHandler(container controls.Container) echo.HandlerFunc {
	return func(c echo.Context) error {
		path := c.Param("*")
		if path == "" {
			return echo.NewHTTPError(http.StatusNotFound, "Path is empty")
		}
		path = sanitizePath(path)
		segments := strings.Split(path, "/")
		item, err := controls.ContainerFollowPath(container, segments)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Error following path '%s': %s", path, err))
		}

		if control, ok := item.(controls.Control); ok {
			var value interface{}
			if err := json.NewDecoder(c.Request().Body).Decode(&value); err != nil {
				return err
			}

			switch v := value.(type) {
			case int:
				control.SetValueString(strconv.Itoa(v))
			case float64:
				floatValue := fmt.Sprintf("%f", v)
				control.SetValueString(floatValue)
			case string:
				control.SetValueString(v)
			case bool:
				control.SetValueString(strconv.FormatBool(v))
			default:
				return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Invalid value type: %T", v))
			}
			return c.JSON(http.StatusOK, control.GetValueString())
		} else {
			return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("Item is not a Control (got '%T')", item))
		}
	}
}

func ColorsGetHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.JSON(http.StatusOK, controls.AllColors)
	}
}
