package main

import (
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(http.ResponseWriter, *http.Request) {
		log.Println("Hello World")
	})

	err := http.ListenAndServe(":9090", nil)
	log.Fatal(err)
}