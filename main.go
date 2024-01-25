package main

import (
    "fmt"
    "net/http"
    "sync"
)

// Defining a queue type and attaching enqueue() and dequeue() methods
type MessageQ struct {
    messages []map[string]string
}
func (self *MessageQ) enqueue(user string, ipAddr string, message string) {
    messageMap := map[string]string{
        "user": user,
        "ipAddr": ipAddr,
        "message": message,
    }
    self.messages = append(self.messages, messageMap)
    if len(self.messages) > 5 {
        self.dequeue()
    }
}
func (self *MessageQ) dequeue() map[string]string {
    if len(self.messages) == 0 {
        return nil
    }
    message := self.messages[0]
    self.messages = self.messages[1:]
    return message
}

// Declaring a global queue for messages
var Messages MessageQ

// Goroutine for "/receive" endpoint
func receiveRoutine(
            wg *sync.WaitGroup, 
            w http.ResponseWriter) {

    // Loop over messages in queue and print to http.ResponseWriter
    for _, msg := range Messages.messages {
        Msg := fmt.Sprintf("Q %s (%s):\n   %s\n", msg["user"], msg["ipAddr"], msg["message"])
        fmt.Fprintf(w, Msg)
        fmt.Printf(Msg)
    }

    // Finish waitgroup
    wg.Done()
}

// Goroutine for "/send" endpoint
func sendRoutine(
            wg *sync.WaitGroup, 
            user string, 
            message string,
            ipAddr string) {

    // Add message to queue
    Messages.enqueue(user, ipAddr, message)

    // Finish waitgroup
    wg.Done()
}

// Handler function for requests to "/receive" endpoint
func handleReceive(w http.ResponseWriter, r *http.Request) {

    // Create a waitgroup
    var wg sync.WaitGroup

    // Pull information from request
    ipAddr := r.RemoteAddr

    fmt.Fprintf(w, "Client: %s\n", ipAddr)
    fmt.Printf("Client: %s\n", ipAddr)
    fmt.Fprintf(w, "RECEIVE REQUEST\n")
    fmt.Printf("RECEIVE REQUEST\n")

    // Add waitgroup for goroutine
    wg.Add(1)
    go receiveRoutine(&wg, w)

    // Wait for goroutine to finish
    wg.Wait()
}

// Handler function for requests to "/send" endpoint
func handleSend(w http.ResponseWriter, r *http.Request) {
    
    // Create a waitgroup
    var wg sync.WaitGroup

    // Pull information from request
    params := r.URL.Query()
    user := params.Get("user")
    message := params.Get("message")
    ipAddr := r.RemoteAddr

    fmt.Fprintf(w, "Client: %s\n", ipAddr)
    fmt.Printf("Client: %s\n", ipAddr)
    fmt.Fprintf(w, "SEND REQUEST\n")
    fmt.Printf("SEND REQUEST\n")

    // Check if user is present
    if user == "" {
        user = "Anonymous"
    }

    // Check if message is present
    if message == "" {
        fmt.Fprintf(w, "No message given\n")
        fmt.Printf("No message given\n")
    } else {

        // Add a waitgroup for each goroutine
        wg.Add(1)
        go sendRoutine(&wg, user, message, ipAddr)
        wg.Wait()
        wg.Add(1)
        go receiveRoutine(&wg, w)

        // Wait for all waitgroups to finish
        wg.Wait()
    }
}

// Main driver function
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
