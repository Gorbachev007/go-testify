package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

func printSuccessMessage(testName string) {
	fmt.Printf("\033[32mТест %s завершен успешно\033[0m\n", testName)
}

func main() {
	http.HandleFunc("/cafe", mainHandle)
	http.ListenAndServe(":8080", nil)
}
