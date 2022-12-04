package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func main() {
	ctx := context.Background()
	done := make(chan struct{})
	mux := http.NewServeMux()
	mux.HandleFunc("/app", ws)
	s := http.Server{
		Addr:              ":6000",
		Handler:           mux,
		TLSConfig:         nil,
		ReadTimeout:       200,
		ReadHeaderTimeout: 200,
		WriteTimeout:      200,
		IdleTimeout:       0,
		MaxHeaderBytes:    1 << 20,
		TLSNextProto:      nil,
		ConnState:         nil,
		ErrorLog:          nil,
		BaseContext:       nil,
		ConnContext:       nil,
	}

	s.ListenAndServe()

	go func() {
		var s string
		fmt.Scanln(&s)
		close(done)
	}()

	<-done
	s.Shutdown(ctx)
}

func ws(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Failed to read message %v", err)
			conn.Close()
			return
		}
		log.Printf(string(msg))

		err = conn.WriteMessage(1, []byte("Server replied"))
		if err != nil {
			log.Printf("Faile to send message %v", err)
			return
		}
	}
}
