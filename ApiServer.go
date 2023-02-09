package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Method + " " + r.URL.String())
		if r.Method == http.MethodGet {
			// Обработка GET запроса
		} else if r.Method == http.MethodPost {
			// Обработка POST запроса
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		newStr := buf.String()
		fmt.Printf("Body request:\n" + newStr + "\n")
	})
	http.ListenAndServe(":8080", nil)
}
