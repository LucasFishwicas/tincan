package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)




// Set upgrader variable to access Upgrader method
var upgrader = websocket.Upgrader{
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
}




// Handler function to handle requests
func handle(w http.ResponseWriter,r *http.Request) {
    // Check if request wants to upgrade to Websocket
    if r.Header.Get("Upgrade") == "websocket" {
        wsHandler(w,r)
    } else {
        httpHandler()
    }
    
}



func wsHandler(w http.ResponseWriter,r *http.Request) {
    // Upgrade to Websocket or print error and return
    conn, err := upgrader.Upgrade(w,r,nil)
    if err != nil {
        log.Println("Failed to Upgrade to Websocket:",err)
        return
    }
    defer conn.Close()

    // Notify server of successful upgrade
    fmt.Println("Successfully Upgraded http to Websocket")

    
    // ERROR ON CLIENT SIDE
    // Receives HTTP 200 OK instead of HTTP 101 SWITCHING PROTOCOL


    // Send a welcome message  --  doesn't appear
    err = conn.WriteMessage(websocket.TextMessage, []byte("Hello from server!"))
    if err != nil {
        fmt.Println("Error sending message:", err)
        return
    }


    // Eternally loop and check messages - causes error on ReadMessage
    /*for {
        messageType,_,err := conn.ReadMessage()
        if err != nil {
            log.Println("Erro reading message:",err)
            return
        }

        if err := conn.WriteMessage(messageType,[]byte("Sent from WriteMessage?")); err != nil {
            log.Println("WriteMessage:",err)
            return
        }
    }*/
}



func httpHandler() {
    // Handle Regular HTTP Request
    fmt.Println("Not a websocket upgrade")
}






func main() {
    // Define a handler function for incoming requests
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w,"Welcome to tincan\nThe single-line CLI chat service\n")
        
        // handle http request or ws upgrade request 
        handle(w,r)
    })
    

    // Start the server on port 8080
    fmt.Println("Server listening on port 8080")
    http.ListenAndServe(":8080", nil)
}

