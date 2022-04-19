// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"
	"net/http"
	"io/fs"
	"embed"
	"log"

	"github.com/jonathongardner/file-share-server/hub"

	"github.com/gorilla/websocket"
)
//go:embed public/*
var uiFolder embed.FS

var addr = flag.String("addr", ":8080", "http service address")

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	flag.Parse()
	h := hub.NewHub()
	go h.Run()
	http.HandleFunc("/candidates", func(w http.ResponseWriter, r *http.Request) {
		ipAddress := readUserIP(r)
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		h.NewClient(conn, ipAddress, "")
	})
	ui, _ := fs.Sub(uiFolder, "public")
	http.Handle("/", http.FileServer(http.FS(ui)))
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func readUserIP(r *http.Request) string {
    IPAddress := r.Header.Get("X-Real-Ip")
    if IPAddress == "" {
        IPAddress = r.Header.Get("X-Forwarded-For")
    }
    if IPAddress == "" {
        IPAddress = r.RemoteAddr
    }
    return IPAddress
}
