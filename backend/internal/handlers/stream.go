package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/amilcar-vasquez/blessed-bites/backend/internal/realtime"
)

type StreamHandler struct {
	Broker *realtime.Broker
}

func NewStreamHandler(broker *realtime.Broker) *StreamHandler {
	return &StreamHandler{Broker: broker}
}

func (h *StreamHandler) OrdersStream(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	ch := h.Broker.Subscribe()
	defer h.Broker.Unsubscribe(ch)

	ping := time.NewTicker(20 * time.Second)
	defer ping.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case <-ping.C:
			fmt.Fprint(w, "event: ping\ndata: {}\n\n")
			flusher.Flush()
		case evt := <-ch:
			payload, _ := json.Marshal(evt.Data)
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", evt.Type, payload)
			flusher.Flush()
		}
	}
}
