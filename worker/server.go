package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/adjust/rmq/v5"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Handler struct {
	conn rmq.Connection
}

func NewHandler(conn rmq.Connection) *Handler {
	return &Handler{
		conn: conn,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	layout := r.FormValue("layout")
	refresh := r.FormValue("refresh")

	queues, err := h.conn.GetOpenQueues()
	if err != nil {
		panic(err)
	}

	stats, err := h.conn.CollectStats(queues)
	if err != nil {
		panic(err)
	}

	fmt.Fprint(w, stats.GetHtml(layout, refresh))
}

type WSHandler struct{}

func NewWSHandler() *WSHandler {
	return &WSHandler{}
}

func (wsHandler *WSHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalln("Http connection upgrade to websockets failed:", err)
	}
	defer conn.Close()

	isClosed := make(chan bool)
	go func(conn *websocket.Conn) {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Printf("websocket readMsg error: %v", err)
			// send over channel on websocket error
			isClosed <- true
		}

	}(conn)

loop:
	for range time.Tick(5 * time.Second) {
		// process sequentially from the queue
		for _, m := range statusQueue {
			msg, err := json.Marshal(m)
			if err != nil {
				log.Printf("bad request: %v", err)
				continue
			}

			//Write the message back to the browser
			if err := conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				log.Println("Unable to write back message")

			}
			statusQueue = statusQueue[1:]
		}

		select {
		case <-isClosed:
			log.Println("received over closed channel, breaking out of the loop")
			break loop
		default:
			log.Println("conn still active")
		}
	}

}

func statsServer(conn rmq.Connection) {
	http.Handle("/overview", NewHandler(conn))

	// for handling orders
	http.Handle("/listen-orders", NewWSHandler())
	http.HandleFunc("/place-order", placeOrder)

	fmt.Printf("Handler listening on http://localhost:3333")
	if err := http.ListenAndServe(":3333", nil); err != nil {
		log.Fatalln(err)
	}
}

func placeOrder(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./worker/place_order.html")
}
