package Controllers

import (
	"TaiBaiSupport/Models"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strconv"
)

var upgrader = websocket.Upgrader{}

func WSHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("good")
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	v := r.URL.Query();
	classroomId, _ := strconv.Atoi(v.Get("classroomId"))
	userId,_ := strconv.Atoi(v.Get("userId"))

	wsEvent := Models.TaibaiUserWsEvent{
		ClassroomId: classroomId,
		UserId:userId,
		Conn:c,
	}

	Models.TaibaiClassroomManagerInstance.PendingWsChan <- wsEvent
}

func EchoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
}