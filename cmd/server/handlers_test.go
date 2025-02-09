package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
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

func getEchoTest(method, path string, body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(path)
	return c, rec
}

func Test_ContainerGetItem(t *testing.T) {
	c, rec := getEchoTest(http.MethodGet, "/container", "")
	c.SetParamNames("name")
	c.SetParamValues("0")

	list := getRootList()

	require.NoError(t, ContainerGetHandler(list)(c))
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "false\n", rec.Body.String())
}

func getRootList() *controls.List {
	list := controls.NewList(1)
	toggle := controls.NewToggle()
	list.SetItem(0, toggle)
	return list
}

func Test_ContainerGetContainer(t *testing.T) {
	c, rec := getEchoTest(http.MethodGet, "/container", "")
	list := getRootList()

	require.NoError(t, ContainerGetHandler(list)(c))
	resp := rec.Result()
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Equal(t, `{"0":false}`+"\n", rec.Body.String())

}
func Test_ContainerGetContainerSlash(t *testing.T) {
	c, rec := getEchoTest(http.MethodGet, "/container/", "")

	list := getRootList()
	require.NoError(t, ContainerGetHandler(list)(c))

	assert.Equal(t, http.StatusOK, rec.Code)

	resp := rec.Result()
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	body := rec.Body.String()
	assert.Equal(t, `["0"]`+"\n", body)
}

func Test_ContainerGetContainerControl(t *testing.T) {
	c, rec := getEchoTest(http.MethodGet, "/container/0", "")
	c.SetParamNames("name")
	c.SetParamValues("0")

	list := getRootList()
	require.NoError(t, ContainerGetHandler(list)(c))
	resp := rec.Result()
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Equal(t, `false`+"\n", rec.Body.String())
}

func Test_ContainerPostSetValue(t *testing.T) {
	c, rec := getEchoTest(http.MethodPost, "/container/0", "true")
	c.SetParamNames("*")
	c.SetParamValues("0")

	err := ContainerPostHandler2(getRootList())(c)
	require.NoError(t, err)
	resp := rec.Result()

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Equal(t, `true`+"\n", rec.Body.String())
}
