package worker

import (
	"fmt"
	"log"
	"net/http"

	"github.com/adjust/rmq/v5"
)

func statsServer(conn rmq.Connection) {
	http.Handle("/overview", NewHandler(conn))

	fmt.Printf("Handler listening on http://localhost:3333/overview\n")
	if err := http.ListenAndServe(":3333", nil); err != nil {
		log.Fatalln(err)
	}
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

	log.Printf("queue stats\n%s", stats)
	fmt.Fprint(w, stats.GetHtml(layout, refresh))
}
