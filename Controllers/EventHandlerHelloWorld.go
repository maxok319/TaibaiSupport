package Controllers

import (
	"fmt"
	"net/http"
)

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}