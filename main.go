package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)


var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
}


func handler(w http.ResponseWriter,r *http.Request) {
    _, err := upgrader.Upgrade(w,r,nil)
    if err != nil {
        log.Println(err)
        return
    }

    fmt.Println("Successfully upgraded http to Websocket")
}


func main() {

    // Define a handler function for incoming requests
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w,"Welcome to tincan\nThe single-line CLI chat service\n")
        fmt.Fprintf(w,r.Header.Get("Sec-WebSocket-Key"))
        handler(w,r)
    })
    

    // Start the server on port 8080
    fmt.Println("Server listening on port 8080")
    http.ListenAndServe(":8080", nil)
}

