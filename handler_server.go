package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func handler(commands Commands) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, r *http.Request) {
		pathSegments := strings.Split(r.URL.Path, "/")
		var command string
		for i, segment := range pathSegments {
			if segment == "api" && i < len(pathSegments)-1 {
				command = pathSegments[i+1]
				break
			}
		}

		var args []Argument
		query := r.URL.Query()

		for k := range query {
			args = append(args, Argument{
				Name:  k,
				Value: query.Get(k),
			})

		}

		fmt.Fprintf(w, "Command: %v\nParams: %v\n", command, args)
	}
}

func handleServer(commands Commands) {
	http.HandleFunc("/api/", handler(commands))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
