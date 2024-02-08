package main

import (
    "fmt"
    "net/http"
	"bufio"
	"log"
	"os"
    //"os/signal"

	"github.com/gorilla/websocket"
	"sync"

    "dev/golang/tincan/models"
)

var (
    // Create Upgrader struct
    upgrader = websocket.Upgrader{
        ReadBufferSize: 1024,
        WriteBufferSize: 1024,
    }
    // Create channel to handle interrupt
    //signalChan = make(chan os.Signal, 1)
    // Create channel for sending server-side to client
    sendChan chan string
    // Create channel for receiving client-side to server
    receiveChan chan string
    // Create channel to handle interrupt
    //signalChan = make(chan os.Signal, 1)
    // Declaring a global queue for messages
    Messages *models.MessageQ
)

func init() {
    // Initialising sendChan
    sendChan = make(chan string)
    // Initialising receiveChan
    receiveChan = make(chan string)
    // Initialising Messages queue of size 5
    Messages = models.CreateQ(5)
}


// NOT USED -- Handler function to handle requests
func handle(w http.ResponseWriter, r *http.Request) {
    // Check if request wants to upgrade to Websocket
    if r.Header.Get("Upgrade") == "websocket" {
        wsHandler(w, r)
    } else {
        fmt.Println("Not a websocket upgrade")
    }
}


// Read message sent from user server-side
func readSend(wg *sync.WaitGroup) {
    scanner := bufio.NewScanner(os.Stdin)
    for scanner.Scan() {
        text := scanner.Text() // Get the current line of text
        if text != "" {
            sendChan <- text
        }
    }
    if err := scanner.Err(); err != nil {
        log.Println("Error reading sent:", err)
    }
    close(sendChan)
    wg.Done()
}


// Read message sent from user client-side
func readReceive(wg *sync.WaitGroup, conn *websocket.Conn, client string) {
    for { 
        messageType, message, err := conn.ReadMessage()
        if err != nil {
            fmt.Println("client://", client,  "left the chat")
            return
        }
        if string(message) != "" {
            receiveChan <- string(message)
        }

        if messageType == websocket.CloseNormalClosure {
            log.Println("Normal Close messageType received")
            break
        }
    }
    close(receiveChan)
    wg.Done()
}


func wsHandler(w http.ResponseWriter, r *http.Request) {
    // Upgrade to Websocket or print error and return
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Failed to Upgrade to Websocket:", err)
        return
    }
    defer conn.Close()

    client := conn.RemoteAddr().String()

    // Notify server of successful upgrade
    fmt.Println("client://", client, "entered the chat")

    // Send a welcome message
    err = conn.WriteMessage(websocket.TextMessage, 
                            []byte("tincan:// You're in! "),
          )
    if err != nil {
        log.Println("Error sending message:", err)
        return
    }

    // Create a waitgroup
    var wg sync.WaitGroup

    // Launch goroutines for reading and writing
    wg.Add(2)
    go readSend(&wg)
    go readReceive(&wg, conn, client)

    // Eternally loop and check channel messages
    for {
        select {  
        case sent := <-sendChan:
            if err := conn.WriteMessage(websocket.TextMessage, 
                                        []byte("-> "+sent),
                      ); err != nil {
                log.Println("Error writing message:", err)
                return
            }
        case received := <-receiveChan:
            fmt.Print("-> ", string(received), "\n") 
        }
    }
}




// Goroutine for "http/receive" endpoint
func httpReceiveRoutine(wgQ *sync.WaitGroup, w http.ResponseWriter) {
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
func httpSendRoutine(
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
func httpHandleReceive(w http.ResponseWriter, r *http.Request) {
    // Create a waitgroup
    var wgQ sync.WaitGroup

    // Pull information from request
    ipAddr := r.RemoteAddr

    fmt.Fprintf(w, "RECEIVE REQUEST: %s\n", ipAddr)
    fmt.Printf("RECEIVE REQUEST: %s\n", ipAddr)

    // Add waitgroup for goroutine
    wgQ.Add(1)
    go httpReceiveRoutine(&wgQ, w)

    // Wait for goroutine to finish
    wgQ.Wait()
}


// Handler function for requests to "/http/send" endpoint
func httpHandleSend(w http.ResponseWriter, r *http.Request) {
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
        go httpSendRoutine(&wgQ, user, message, ipAddr)
        wgQ.Wait()
        //wgQ.Add(1)
        //go httpReceiveRoutine(&wgQ, w)

        // Wait for all waitgroups to finish
        //wgQ.Wait()
    }
}




// Main driver function
func main() {
    // Define a handler function for Home
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w,"Welcome to tincan://\nThe single-line CLI chat service\n")
    })

    // Handler function for establishing websocket
    http.HandleFunc("/ws", wsHandler)

    // Handler function for sending URL parameter messages. "/send?message=..."
    http.HandleFunc("/http/send", httpHandleSend)

    // Handler function for receiving URL parameter messages.
    http.HandleFunc("/http/receive", httpHandleReceive)

    // Start the server on port 8080
    fmt.Println("Server listening on port 8080")
    http.ListenAndServe(":8080", nil)

    // Wait for goroutine to finish
    //wg.Wait()
    fmt.Println("All goroutines finished")
}
