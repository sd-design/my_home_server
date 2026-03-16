package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)
	mux.HandleFunc("/api/disk", listDirectoryHandler)
	mux.HandleFunc("/tech/my_ip", getClientIp)
	mux.HandleFunc("/filemanager", showfileManager)

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	mux.Handle("/static/", http.StripPrefix("/static", fileServer))

	log.Println("Starting server on :4200")
	err := http.ListenAndServe(":4200", mux)
	log.Fatal(err)

}
