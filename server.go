package main

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/gorilla/websocket"
)

//go:embed dist
var distFiles embed.FS

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// Allow all origins — appropriate for a trusted LAN game with no authentication.
	CheckOrigin: func(r *http.Request) bool { return true },
}

func newServer(hub *Hub) http.Handler {
	mux := http.NewServeMux()

	// Strip the "dist/" prefix so index.html is served at "/".
	sub, err := fs.Sub(distFiles, "dist")
	if err != nil {
		panic(err)
	}
	mux.Handle("/", http.FileServer(http.FS(sub)))

	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWS(hub, upgrader, w, r)
	})

	return mux
}
