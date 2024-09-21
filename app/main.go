package main

import (
	"fmt"
	"net/http"
	"text/template"
)

// обработка html
func homePage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("../web/index.html")
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	tmp.Execute(w, nil)
}

func handleRequest() {
	// обработка css и js
	http.Handle("/web/", http.StripPrefix("/web/", http.FileServer(http.Dir("../web/"))))
	// обработка изображений
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("../assets/"))))
	// Отображение страниц
	http.HandleFunc("/", homePage)
	http.ListenAndServe(":8080", nil)
}

func main() {
	handleRequest()
}
