package main

import (
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

func Test_ContainerPostHandler(t *testing.T) {
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
			c, rec := getResponseRecorder(http.MethodPost, "/api/v2/root/"+tt.path, "")
			c.SetParamNames("*")
			c.SetParamValues(tt.path)
			c.Request().Body = io.NopCloser(strings.NewReader(tt.value))

			// Execute handler
			err := ContainerPostHandler(rootContainer)(c)

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

func getResponseRecorder(method, path string, body string) (echo.Context, *httptest.ResponseRecorder) {
	e := echo.New()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath(path)
	return c, rec
}

func getRootList() *controls.List {
	list := controls.NewList(1)
	toggle := controls.NewToggle()
	list.SetItem(0, toggle)
	return list
}

func Test_ContainerPostSetValue(t *testing.T) {
	c, rec := getResponseRecorder(http.MethodPost, "/container/0", "true")
	c.SetParamNames("*")
	c.SetParamValues("0")

	err := ContainerPostHandler(getRootList())(c)
	require.NoError(t, err)
	resp := rec.Result()

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Equal(t, `true`+"\n", rec.Body.String())
}

func getTestContainer() controls.Container {
	rootContainer := controls.NewItemMap()
	rootContainer["dial1"] = controls.NewNumericDial()
	container1 := controls.NewItemMap()
	rootContainer["container1"] = container1
	dial2 := controls.NewNumericDial()
	container1["dial2"] = dial2
	dial2.SetValue(42)

	return rootContainer
}

func newGetResponseRecorder(path string) (echo.Context, *httptest.ResponseRecorder) {
	ctx, rec := getResponseRecorder(http.MethodGet, "/*", "")
	ctx.SetParamNames("*")
	ctx.SetParamValues(path)
	return ctx, rec
}

func Test_ContainerGetDial(t *testing.T) {
	handler := ContainerGetHandler(getTestContainer())

	ctx, recorder := newGetResponseRecorder("/container1/dial2")

	require.NoError(t, handler(ctx))
	resp := recorder.Result()

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Equal(t, `42`+"\n", recorder.Body.String())
}

func Test_ContainerGetContainer(t *testing.T) {
	handler := ContainerGetHandler(getTestContainer())

	ctx, recorder := newGetResponseRecorder("/container1")

	require.NoError(t, handler(ctx))
	resp := recorder.Result()

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Contains(t, recorder.Body.String(), `"dial2"`)
	assert.Contains(t, recorder.Body.String(), `42`)
}

func Test_ContainerGetContainerSlash(t *testing.T) {
	handler := ContainerGetHandler(getTestContainer())

	ctx, recorder := newGetResponseRecorder("/container1/")

	require.NoError(t, handler(ctx))
	resp := recorder.Result()

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Contains(t, recorder.Body.String(), `"dial2"`)
	assert.Contains(t, recorder.Body.String(), `[`)
}
