package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var cafeList = map[string][]string{
	"moscow": {"Мир кофе", "Сладкоежка", "Кофе и завтраки", "Сытый студент"},
}

func mainHandle(w http.ResponseWriter, req *http.Request) {
	countStr := req.URL.Query().Get("count")
	if countStr == "" {
		http.Error(w, "count missing", http.StatusBadRequest)
		return
	}

	count, err := strconv.Atoi(countStr)
	if err != nil {
		http.Error(w, "wrong count value", http.StatusBadRequest)
		return
	}

	city := req.URL.Query().Get("city")
	if city == "" {
		http.Error(w, "city missing", http.StatusBadRequest)
		return
	}

	cafe, ok := cafeList[city]
	if !ok {
		http.Error(w, "wrong city value", http.StatusBadRequest)
		return
	}

	if count > len(cafe) {
		count = len(cafe)
	}

	answer := strings.Join(cafe[:count], ",")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(answer))
}

func TestMainHandlerCorrectRequest(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=moscow", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.NotEmpty(t, responseRecorder.Body.String())
}

func TestMainHandlerWrongCity(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2&city=invalidcity", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(t, "wrong city value\n", responseRecorder.Body.String())
}

func TestMainHandlerCountMoreThanTotal(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=10&city=moscow", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	expectedAnswer := strings.Join(cafeList["moscow"], ",")

	assert.Equal(t, http.StatusOK, responseRecorder.Code)
	assert.Equal(t, expectedAnswer, responseRecorder.Body.String())
}

func TestMainHandlerMissingCount(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?city=moscow", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(t, "count missing\n", responseRecorder.Body.String())
}

func TestMainHandlerInvalidCount(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=invalid&city=moscow", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(t, "wrong count value\n", responseRecorder.Body.String())
}

func TestMainHandlerMissingCity(t *testing.T) {
	req := httptest.NewRequest("GET", "/cafe?count=2", nil)
	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(mainHandle)
	handler.ServeHTTP(responseRecorder, req)

	assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
	assert.Equal(t, "city missing\n", responseRecorder.Body.String())
}

func printSuccessMessage(testName string) {
	fmt.Printf("\033[32mТест %s завершен успешно\033[0m\n", testName)
}

func main() {
	t := &testing.T{}

	TestMainHandlerCorrectRequest(t)
	printSuccessMessage("Правильный запрос")

	TestMainHandlerWrongCity(t)
	printSuccessMessage("Неверный город")

	TestMainHandlerCountMoreThanTotal(t)
	printSuccessMessage("Количество больше общего")

	TestMainHandlerMissingCount(t)
	printSuccessMessage("Отсутствует count")

	TestMainHandlerInvalidCount(t)
	printSuccessMessage("Неверное значение count")

	TestMainHandlerMissingCity(t)
	printSuccessMessage("Отсутствует город")

	http.HandleFunc("/cafe", mainHandle)
	http.ListenAndServe(":8080", nil)
}
