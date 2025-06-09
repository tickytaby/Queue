package main

import (
	"fmt"
	"log"
	"net/http"
	
	"github.com/gorilla/websocket"
)

type PurchaseRequest struct{
	Msg string
	Conn *websocket.Conn
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool { return true },
}

func handleWebSocket(w http.ResponseWriter, r *http.Request, queue chan PurchaseRequest) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket: ", err)
		return
	}
	log.Println("New WS Connection Established")
	
	go handleClient(conn, queue)
}

func handleClient(conn *websocket.Conn, queue chan PurchaseRequest) {
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("Read error: ", err)
			break
		}

		msg := string(message)
		
		select {
		case queue <- PurchaseRequest{
			Msg: msg,
			Conn: conn,
		}:
			conn.WriteMessage(websocket.TextMessage, []byte("Successfully Enqueued"))
		default:
			conn.WriteMessage(websocket.TextMessage, []byte("Queue is full"))
		}


		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error writing: ", err)
			break
		}
	}
}

func QueueMonitor(queue <-chan PurchaseRequest) {
	for req := range queue {
		
		fmt.Println("Processing: ", req.Msg)

		response := fmt.Sprintf("You bought: %s", req.Msg)

		if req.Conn != nil {
			err := req.Conn.WriteMessage(websocket.TextMessage, []byte(response))
			if err != nil {
				fmt.Println("Failed to write to WS: ", err)
			}
		}
	}
}

func main() {
	queue := make(chan PurchaseRequest, 100)
	go QueueMonitor(queue)
	fmt.Println("Welcome to the flash queue!")
	http.HandleFunc("/", func(w http.ResponseWriter, r* http.Request) {
		handleWebSocket(w, r, queue)
	})
	log.Println("WebSocket server started on :8080")
	if err := http.ListenAndServe("0.0.0.0:8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
