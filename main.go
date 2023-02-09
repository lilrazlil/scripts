package main

import (
   "fmt"
   "net/http"
   "io/ioutil"
)

func main() {
   http.HandleFunc("/serverstatus", post)
   fmt.Println("Listening on port 8080")
   http.ListenAndServe(":8080", nil)
}

func post(w http.ResponseWriter, r *http.Request) {
   defer r.Body.Close()
   body, err := ioutil.ReadAll(r.Body)
   if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
   }
   if r.Method == http.MethodGet {
		// Обработка GET запроса
		if string(body) == "test1" {
			fmt.Fprintf(w, "GET")
		}

   } else if r.Method == http.MethodPost {
		// Обработка POST запроса
		if string(body) == "test2" {
			fmt.Fprintf(w, "POST")
		}
   }
   // do something here with body
   fmt.Println(string(body))
   // ... run script here ...
   fmt.Fprintf(w, "Success")
}