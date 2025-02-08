package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/mholzen/play-go/controls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TestContainer struct {
	items map[string]controls.Item
}

func NewTestContainer() *TestContainer {
	return &TestContainer{
		items: make(map[string]controls.Item),
	}
}

func (c *TestContainer) GetItem(name string) (controls.Item, error) {
	if item, ok := c.items[name]; ok {
		return item, nil
	}
	return nil, fmt.Errorf("item not found: '%s'", name)
}

func (c *TestContainer) Items() map[string]controls.Item {
	return c.items
}

func Test_ContainerPostHandler2(t *testing.T) {
	// Create a nested container structure for testing
	innerDial := controls.NewNumericDial()
	innerContainer := controls.NewItemMap()
	innerContainer["dial"] = innerDial

	rootContainer := controls.NewItemMap()
	rootContainer["inner"] = innerContainer

	// Test cases
	tests := []struct {
		name           string
		path           string
		value          string
		expectedStatus int
		expectedValue  byte
	}{
		{
			name:           "single level path",
			path:           "inner",
			value:          "42",
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "nested path",
			path:           "inner/dial",
			value:          "123",
			expectedStatus: http.StatusOK,
			expectedValue:  123,
		},
		{
			name:           "invalid path",
			path:           "nonexistent/dial",
			value:          "42",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "path through non-container",
			path:           "inner/dial/invalid",
			value:          "42",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup request
			c, rec := getEchoTest(http.MethodPost, "/api/v2/root/"+tt.path, "")
			c.SetParamNames("*")
			c.SetParamValues(tt.path)
			c.Request().Body = io.NopCloser(strings.NewReader(tt.value))

			// Execute handler
			err := ContainerPostHandler2(rootContainer)(c)

			// Check response
			if tt.expectedStatus != http.StatusOK {
				assert.Error(t, err)
				he, ok := err.(*echo.HTTPError)
				require.True(t, ok, "expected error to be an HTTPError, got %s", err)
				assert.Equal(t, tt.expectedStatus, he.Code, he.Message)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedStatus, rec.Code)

				// For successful cases, verify the value was set correctly
				// switch tt.path {
				// case "inner/dial":
				// 	assert.Equal(t, tt.expectedValue, innerDial.Value)
				// case "inner":
				// 	container, ok := rootContainer["inner"].(controls.Container)
				// 	require.True(t, ok)
				// 	dial, ok := container.GetItem("dial")
				// 	require.True(t, ok)
				// 	assert.Equal(t, tt.expectedValue, dial.Value)
				// }
			}
		})
	}
}
