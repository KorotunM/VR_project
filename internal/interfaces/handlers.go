package interfaces

import (
	"fmt"
	"html/template"
	"net/http"
)

func HomePage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("../web/index.html")
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	tmp.Execute(w, nil)
}

func AdminPage(w http.ResponseWriter, r *http.Request) {
	tmp, err := template.ParseFiles("../web/admin.html")
	if err != nil {
		fmt.Fprintf(w, "Error: %v", err)
		return
	}
	tmp.Execute(w, nil)
}
