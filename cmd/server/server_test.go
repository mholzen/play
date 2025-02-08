package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/mholzen/play-go/controls"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
	c.SetParamNames("name", "value")
	c.SetParamValues("0", "true")

	list := getRootList()
	item0, err := list.GetItem("0")
	require.NoError(t, err)

	control := item0.(controls.Control)
	control.SetValueString("true")
	require.NoError(t, ContainerPostHandler(list)(c))
	resp := rec.Result()

	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	assert.Equal(t, `true`+"\n", rec.Body.String())
}
