package main

import (
    "fmt"
    "net/http"

    "dev/golang/tincan/handlers"
)


// Main driver function
func main() {
    // Define a handler function for Home
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w,"Welcome to tincan://\nThe single-line CLI chat service\n")
    })

    // Handler function for establishing websocket
    http.HandleFunc("/ws", handlers.WsHandler)

    // Handler function for sending URL parameter messages. "/send?message=..."
    http.HandleFunc("/http/send", handlers.HttpHandleSend)

    // Handler function for receiving URL parameter messages.
    http.HandleFunc("/http/receive", handlers.HttpHandleReceive)

    // Start the server on port 8080
    fmt.Println("Server listening on port 8080")
    http.ListenAndServe(":8080", nil)

    // Wait for goroutine to finish
    //wg.Wait()
    //fmt.Println("All goroutines finished")
}
