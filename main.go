package main

import (
	"fmt"
	"net/http"
	"sync"
)

// Create a channel and waitgroup
var messageChan = make(chan string, 100) 
var wg sync.WaitGroup

func receiveRoutine(w http.ResponseWriter) {

    // Loop over messages in channel and print to http.ResponseWriter
    for message := range messageChan {
        fmt.Fprintf(w, message)
    }

    // Finish waitgroup
    wg.Done()
}

func sendRoutine(user string, message string) {

    // Add message to channel
    messageChan <- fmt.Sprintf("%s:\n   %s\n", user, message)

    // Close channel and finish waitgroup
    close(messageChan)
    wg.Done()
}

func handleReceive(w http.ResponseWriter, r *http.Request) {

    fmt.Fprintf(w, "RECEIVE REQUEST\n")
    fmt.Printf("RECEIVE REQUEST\n")

    wg.Add(1)
    go receiveRoutine(w)
    wg.Wait()
}

func handleSend(w http.ResponseWriter, r *http.Request) {
    
    fmt.Fprintf(w, "SEND REQUEST\n")
    fmt.Printf("SEND REQUEST\n")


    // Pull message from URL parameters (e.g., ?message=...)
    params := r.URL.Query()
    user := params.Get("user")
    message := params.Get("message")

    // Check if user is present
    if user == "" {
        user = "Anonymous"
    }

    // Check if message is present
    if message == "" {
        fmt.Fprintf(w, "No message given\n")
    } else {
    
        // Add a waitgroup for each goroutine
        wg.Add(1)
        go sendRoutine(user, message)
        wg.Add(1)
        go receiveRoutine(w)

        // Wait for all waitgroups to finish
        wg.Wait()
    }
}

func main() {

    // Define a handler function for Home
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w,"Welcome to tincan\nThe single-line CLI chat service")
    })

    // Handler function for sending URL parameter messages. "/send?message=..."
    http.HandleFunc("/send", handleSend)

    // Handler function for receiving URL parameter messages.
    http.HandleFunc("/receive", handleReceive)

    // Start the server on port 8080
    fmt.Println("Server listening on port 8080")
    http.ListenAndServe(":8080", nil)
}
