package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
)

type PageData struct {
	PageTitle string
	Message   string
	Items     []string
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	Message := "Your IP-address: " + r.RemoteAddr
	data := PageData{
		PageTitle: "My home server – Cloaca v.1.0",
		Message:   Message,
	}

	files := []string{
		"./ui/html/basic.body.tmpl",
		"./ui/html/header.tmpl",
		"./ui/html/footer.tmpl",
	}

	w.Header().Add("Engine", "Cloaca v.1.0")
	ts, err := template.ParseFiles(files...)

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
		return
	}
	err = ts.Execute(w, data)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal Server Error", 500)
	}

}

func showSnippet(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "Display a specific snippet with ID %d...", id)
}
func createSnippet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", 405)
		return
	}
	w.Write([]byte("Create a new snippet..."))
}

func getClientIp(w http.ResponseWriter, r *http.Request) {
	ipClient := r.RemoteAddr
	fmt.Fprintf(w, "Your Ip-adress %s", ipClient)
}
