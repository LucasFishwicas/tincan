package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)




var (
    // Set upgrader variable to access Upgrader method
    upgrader = websocket.Upgrader{
        ReadBufferSize: 1024,
        WriteBufferSize: 1024,
    }
    // Create channel for sending server-side input to client
    inputChan = make(chan string)
)




// Handler function to handle requests
func handle(w http.ResponseWriter,r *http.Request) {
    // Check if request wants to upgrade to Websocket
    if r.Header.Get("Upgrade") == "websocket" {
        wsHandler(w,r)
    } else {
        httpHandler()
    }
    
}


// Read input from user server-side
func readInput() {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        text := scanner.Text() // Get the current line of text
        if text != "" {
            inputChan <- text
        }
    }
    if err := scanner.Err(); err != nil {
        log.Println("Error reading input:", err)
    }
    close(inputChan)
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

    // Send a welcome message  --  doesn't appear
    err = conn.WriteMessage(websocket.TextMessage, []byte("Hello from server!"))
    if err != nil {
        fmt.Println("Error sending message:", err)
        return
    }


    // Eternally loop and check messages - causes error on ReadMessage
    for {
        select {
            case input := <-inputChan: 
            if err := conn.WriteMessage(websocket.TextMessage,[]byte(input)); err != nil {
                log.Println("Error writing message:",err)
                return
            }
        default:
            messageType,_,err := conn.ReadMessage()
            if err != nil {
                log.Println("Error reading message:",err)
                return
            }
            log.Println(messageType)
        }
    }
}



func httpHandler() {
    // Handle Regular HTTP Request
    fmt.Println("Not a websocket upgrade")
}






func main() {
    // Define a handler function for incoming requests
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w,"Welcome to tincan\nThe single-line CLI chat service\n") 
    })
    go http.HandleFunc("/ws",wsHandler)
    
    go readInput() 

    // Start the server on port 8080
    fmt.Println("Server listening on port 8080")
    http.ListenAndServe(":8080", nil)
}

