package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/command", commandHandler)
	http.HandleFunc("/api/getstatus", getStatus)
	fmt.Println("Brook-pf server starting")
	http.ListenAndServe(":8000", nil)
}
