package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("Starting HTTPS Server")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello Encrypted World")
	})

	if err := http.ListenAndServeTLS(":8080", "certificate.pem", "key.pem", nil); err != nil {
		log.Fatalln(err)
	}
}
