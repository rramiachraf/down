package main

import (
	"fmt"
	log "log/slog"
	"net"
	"os"
	"strconv"
)

func main() {
	hl := log.NewTextHandler(os.Stdout, &log.HandlerOptions{})
	l := log.New(hl)

	var port int16 = 3000
	envPORT := os.Getenv("PORT")
	if p, err := strconv.Atoi(envPORT); err == nil {
		port = int16(p)
	}

	addr := fmt.Sprintf(":%d", port)
	ls, err := net.Listen("tcp", addr)
	if err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}

	msg := fmt.Sprintf("server is listening at :%d", port)
	l.Info(msg)

	if err := startServer(ls, l); err != nil {
		l.Error(err.Error())
		os.Exit(1)
	}
}
