package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/democ", demoCommandHandler)
	http.HandleFunc("/command", commandHandler)
	http.ListenAndServe(":8000", nil)
	fmt.Println("Brook-pf server starting")
}
