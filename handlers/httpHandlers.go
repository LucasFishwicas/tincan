package handlers
// HTTP handler functions

import (
    "fmt"
    "net/http"
    "sync"

    "dev/golang/tincan/models"
)

var (
    Messages *models.MessageQ
)

func init() {
    Messages = models.CreateQ(5)
}


// Goroutine for "http/receive" endpoint
func HttpReceiveRoutine(wgQ *sync.WaitGroup, w http.ResponseWriter) {
    // Loop over messages in queue and print to http.ResponseWriter
    for i := 0; i < Messages.Length; i++ {
        index := (Messages.Head+i) % Messages.Capacity
        fmt.Printf("length: %d, i: %d, index: %d\n", Messages.Length, i, index)
        msg := Messages.Messages[index]
        if msg["message"] == "" {
            continue
        }
        Msg := fmt.Sprintf("[%d] %s (%s):\n   %s\n", index,
                                                     msg["user"], 
                                                     msg["ipAddr"], 
                                                     msg["message"],
        )
        fmt.Fprintf(w, Msg)
        fmt.Printf(Msg)
    }

    // Finish waitgroup
    wgQ.Done()
}


// Goroutine for "http/send" endpoint
func HttpSendRoutine(
            wgQ *sync.WaitGroup, 
            user string, 
            message string,
            ipAddr string) {
    // Add message to queue
    Messages.Enqueue(user, ipAddr, message)

    // Finish waitgroup
    wgQ.Done()
}


// Handler function for requests to "/http/receive" endpoint
func HttpHandleReceive(w http.ResponseWriter, r *http.Request) {
    // Create a waitgroup
    var wgQ sync.WaitGroup

    // Pull information from request
    ipAddr := r.RemoteAddr

    fmt.Fprintf(w, "RECEIVE REQUEST: %s\n", ipAddr)
    fmt.Printf("RECEIVE REQUEST: %s\n", ipAddr)

    // Add waitgroup for goroutine
    wgQ.Add(1)
    go HttpReceiveRoutine(&wgQ, w)

    // Wait for goroutine to finish
    wgQ.Wait()
}


// Handler function for requests to "/http/send" endpoint
func HttpHandleSend(w http.ResponseWriter, r *http.Request) {
    // Create a waitgroup
    var wgQ sync.WaitGroup

    // Pull information from request
    params := r.URL.Query()
    user := params.Get("user")
    message := params.Get("message")
    ipAddr := r.RemoteAddr

    fmt.Fprintf(w, "SEND REQUEST: %s\n", ipAddr)
    fmt.Printf("SEND REQUEST: %s\n", ipAddr)

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
        wgQ.Add(1)
        go HttpSendRoutine(&wgQ, user, message, ipAddr)
        wgQ.Wait()
        //wgQ.Add(1)
        //go HttpReceiveRoutine(&wgQ, w)

        // Wait for all waitgroups to finish
        //wgQ.Wait()
    }
}

