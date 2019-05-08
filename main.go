package main

import (
	"TaiBaiSupport/Controllers"
	"log"
	"net/http"
)

func main() {
	log.Println("start taibai support")

	http.HandleFunc("/ws", Controllers.WSHandler)
	http.HandleFunc("/echo", Controllers.EchoHandler)

	http.ListenAndServe("0.0.0.0:8888", nil)

	log.Println("stop taibai support")
}
