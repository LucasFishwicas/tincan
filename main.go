package main

import (
    "fmt"
    "net/http"
    "sync"
)

func receiveRoutine(wg *sync.WaitGroup, messageChan chan string, w http.ResponseWriter) {

    // Loop over messages in channel and print to http.ResponseWriter
    for message := range messageChan {
        fmt.Fprintf(w, message)
    }

    // Finish waitgroup
    wg.Done()
}

func sendRoutine(wg *sync.WaitGroup, messageChan chan string, user string, message string) {

    // Add message to channel
    messageChan <- fmt.Sprintf("%s:\n   %s\n", user, message)

    // Close channel and finish waitgroup
    close(messageChan)
    wg.Done()
}

func handleReceive(w http.ResponseWriter, r *http.Request) {

    // Create a channel and waitgroup
    var messageChan = make(chan string, 100) 
    var wg sync.WaitGroup

    fmt.Fprintf(w, "RECEIVE REQUEST\n")
    fmt.Printf("RECEIVE REQUEST\n")

    messageChan <- fmt.Sprintf("CANNOT ACCESS MESSAGES AT THIS TIME")

    wg.Add(1)
    close(messageChan)
    go receiveRoutine(&wg, messageChan, w)
    wg.Wait()
}

func handleSend(w http.ResponseWriter, r *http.Request) {
    
    fmt.Fprintf(w, "SEND REQUEST\n")
    fmt.Printf("SEND REQUEST\n")

    fmt.Fprintf(w, "Client: %s\n", r.RemoteAddr)
    fmt.Printf("Client: %s\n", r.RemoteAddr)

    // Create a channel and waitgroup
    var messageChan = make(chan string, 100) 
    var wg sync.WaitGroup


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
        go sendRoutine(&wg, messageChan, user, message)
        wg.Add(1)
        go receiveRoutine(&wg, messageChan, w)

        // Wait for all waitgroups to finish
        wg.Wait()
    }
}

func main() {

    // Define a handler function for Home
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w,"Welcome to tincan\nThe single-line CLI chat service\n")
    })

    // Handler function for sending URL parameter messages. "/send?message=..."
    http.HandleFunc("/send", handleSend)

    // Handler function for receiving URL parameter messages.
    http.HandleFunc("/receive", handleReceive)

    // Start the server on port 8080
    fmt.Println("Server listening on port 8080")
    http.ListenAndServe(":8080", nil)
}
