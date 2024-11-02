package handlers

import (
	"fmt"
	"net/http"
)

func Greet(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintln(w, "Hello, web!")
}
