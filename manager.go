package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)


var (
    websocketUpgrader = websocket.Upgrader{
        ReadBufferSize: 1024,
        WriteBufferSize: 1024,
    }
)

type Manager struct {
}

func newManager() *Manager {
    return &Manager{}
}

func (m *Manager) serveWS(w http.ResponseWriter, r *http.Request) {
    log.Println("New connection")
    // Begin by upgrading the HTTP request
    conn, err := websocketUpgrader.Upgrade(w,r,nil)
    if err != nil {
        log.Println(err)
        return
    }
    // We wont do anything yet so close the connection again
    conn.Close()
}



