package main

import (
	"TaibaiSupport/Controllers"
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	log.Println("start taibai support")

	http.HandleFunc("/ws", Controllers.HandleEventPendingWS)
	http.HandleFunc("/echo", Controllers.EchoHandler)

	http.ListenAndServe("0.0.0.0:8888", nil)

	log.Println("stop taibai support")
}
