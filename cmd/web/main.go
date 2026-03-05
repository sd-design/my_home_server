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
	mux.HandleFunc("/api/disk", readFolder)
	mux.HandleFunc("/tech/my_ip", getClientIp)
	log.Println("Starting server on :4200")
	err := http.ListenAndServe(":4200", mux)
	log.Fatal(err)

}
