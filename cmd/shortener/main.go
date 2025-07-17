package main

import (
	"fmt"
	"net/http"
)

var nameUrl string

func postHandler(w http.ResponseWriter, r *http.Request) {

	nameUrl = r.URL.Path

	fmt.Fprintln(w, nameUrl)

}

func getHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, nameUrl)

}

func main() {
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/post", postHandler)

	http.ListenAndServe("localhost:8080", nil)

}
