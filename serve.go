package main

import (
	"bytes"
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/rramiachraf/observe/monitor"
	"github.com/rramiachraf/observe/templates"
)

//go:embed static
var static embed.FS

func startServer(listener net.Listener, lg *slog.Logger) error {
	mux := http.NewServeMux()

	mt := monitor.NewMonitor()

	go func() {
		for {
			err := mt.Update()
			if err != nil {
				msg := fmt.Sprintf("unable to update data, %s", err)
				lg.Error(msg)
			}

			time.Sleep(3 * time.Second)
		}
	}()

	mux.Handle("GET /", mainPageHandler{mt, lg})
	mux.Handle("GET /ws", wsHandler{mt, lg})
	mux.Handle("GET /api", apiHandler{mt, lg})

	mux.Handle("GET /static/", http.FileServerFS(static))

	server := &http.Server{
		ReadTimeout:  time.Second * 20,
		WriteTimeout: time.Second * 20,
		Handler:      mux,
	}

	return server.Serve(listener)
}

type mainPageHandler struct {
	monitor *monitor.Monitor
	l       *slog.Logger
}

func (h mainPageHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := templates.MainPage(h.monitor).Render(context.Background(), w)
	if err != nil {
		h.l.Error(err.Error())
		w.WriteHeader(500)
	}
}

var upgrade = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type wsHandler struct {
	monitor *monitor.Monitor
	l       *slog.Logger
}

func (h wsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrade.Upgrade(w, r, nil)
	if err != nil {
		h.l.Error(err.Error())
		return
	}

	defer conn.Close()

	for {
		buf := bytes.NewBuffer(nil)
		if err := templates.Monitor(h.monitor).Render(context.Background(), buf); err != nil {
			h.l.Error(err.Error())
			break
		}

		if err := conn.WriteMessage(1, buf.Bytes()); err != nil {
			h.l.Error(err.Error())
			break
		}

		time.Sleep(3 * time.Second)
	}
}

type apiHandler struct {
	monitor *monitor.Monitor
	l       *slog.Logger
}

func (h apiHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	enc := json.NewEncoder(w)
	w.Header().Set("content-type", "application/json")

	if err := enc.Encode(h.monitor); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		h.l.Error(err.Error())
	}
}
