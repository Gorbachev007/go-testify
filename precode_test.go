package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainHandlerCorrectRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=moscow", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.NotEmpty(t, responseRecorder.Body.String())
	printSuccessMessage("Правильный запрос")
}

func TestMainHandlerWrongCity(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=invalidcity", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(t, "wrong city value\n", responseRecorder.Body.String())
	printSuccessMessage("Неверный город")
}

func TestMainHandlerCountMoreThanTotal(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=10&city=moscow", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	expectedAnswer := strings.Join(cafeList["moscow"], ",")

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, expectedAnswer, responseRecorder.Body.String())
	printSuccessMessage("Количество больше общего")
}

func TestMainHandlerMissingCount(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?city=moscow", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(t, "count missing\n", responseRecorder.Body.String())
	printSuccessMessage("Отсутствует count")
}

func TestMainHandlerInvalidCount(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=invalid&city=moscow", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(t, "wrong count value\n", responseRecorder.Body.String())
	printSuccessMessage("Неверное значение count")
}

func TestMainHandlerMissingCity(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(t, "city missing\n", responseRecorder.Body.String())
	printSuccessMessage("Отсутствует город")
}
