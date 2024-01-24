package main

import (
	"fmt"
	"net/http"
	"sync"
)

func messageReceive(wg *sync.WaitGroup, messageChan chan string, message string) {

    // Add message to channel
    messageChan <- fmt.Sprintf("Message: %s\n", message)

    // Close channel and finish waitgroup
    close(messageChan)
    wg.Done()
}

func messageDisplay(wg *sync.WaitGroup, messageChan chan string, w http.ResponseWriter) {

    // Loop over messages in channel and print to http.ResponseWriter
    for message := range messageChan {
        fmt.Fprintf(w, message)
    }

    // Finish waitgroup
    wg.Done()
}

func handleMessage(w http.ResponseWriter, r *http.Request) {

    // Create a channel and waitgroup
    var messageChan = make(chan string) 
    var wg sync.WaitGroup

    // Pull message from URL parameters (e.g., ?message=...)
    params := r.URL.Query()
    message := params.Get("message")

    // Check if message is present
    if message == "" {
        fmt.Fprintf(w, "No message given\n")
    } else {
    
        // Add a waitgroup for each goroutine
        wg.Add(1)
        go messageReceive(&wg, messageChan, message)
        wg.Add(1)
        go messageDisplay(&wg, messageChan, w)

        // Wait for all waitgroups to finish
        wg.Wait()
    }
}

func main() {
    // Define a handler function for incoming requests
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w,"Welcome to tincan\nThe single-line CLI chat service")
    })

    // Handler function for URL parameter messages. "/chat?message=..."
    http.HandleFunc("/chat", handleMessage)

    // Start the server on port 8080
    fmt.Println("Server listening on port 8080")
    http.ListenAndServe(":8080", nil)
}
