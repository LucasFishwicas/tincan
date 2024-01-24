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



var (
    // Set upgrader variable to access Upgrader method
    upgrader = websocket.Upgrader{
        ReadBufferSize: 1024,
        WriteBufferSize: 1024,
    }

    // Create channel to handle interrupt
    //signalChan = make(chan os.Signal, 1)
    // Create channel for sending server-side to client
    sendChan = make(chan string)
    // Create channel for receiving client-side to server
    receiveChan = make(chan string)
    // Create a waitgroup
    wg sync.WaitGroup
)




// NOT USED -- Handler function to handle requests
func handle(w http.ResponseWriter,r *http.Request) {
    // Check if request wants to upgrade to Websocket
    if r.Header.Get("Upgrade") == "websocket" {
        wsHandler(w,r)
    } else {
        httpHandler()
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
func readReceive(wg *sync.WaitGroup, conn *websocket.Conn) {
    for { 
        messageType,message,err := conn.ReadMessage()
        if err != nil {
            fmt.Println("client:// left the chat")
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






func wsHandler(w http.ResponseWriter,r *http.Request) {
    // Upgrade to Websocket or print error and return
    conn, err := upgrader.Upgrade(w,r,nil)
    if err != nil {
        log.Println("Failed to Upgrade to Websocket:",err)
        return
    }
    defer conn.Close()


    client := conn.RemoteAddr()

    // Notify server of successful upgrade
    fmt.Println("client://",client," entered the chat")

    // Send a welcome message
    err = conn.WriteMessage(websocket.TextMessage, []byte("tincan:// You're in! "))
    if err != nil {
        log.Println("Error sending message:", err)
        return
    }

    // Launch goroutines for reading and writing
    wg.Add(2)
    go readSend(&wg)
    go readReceive(&wg, conn)


    // Eternally loop and check channel messages
    for {
        select {  
        case sent := <-sendChan:
            if err := conn.WriteMessage(websocket.TextMessage,[]byte("-> "+sent)); err != nil {
                log.Println("Error writing message:",err)
                return
            }
        case received := <-receiveChan:
            fmt.Print("-> ",string(received)) 
        }
    }
}


// Reserved to handle http requests
func httpHandler() {
    // Handle Regular HTTP Request
    // Sams code can go here...
    fmt.Println("Not a websocket upgrade")
}




func main() {

    // Define a handler function for incoming requests
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w,"Welcome to tincan\nThe single-line CLI chat service\n") 
    })

    http.HandleFunc("/ws",wsHandler)
    
    // Handler function for URL parameter messages. "/chat?message=..."
    http.HandleFunc("/chat", handleMessage)

    // Start the server on port 8080
    fmt.Println("Server listening on port 8080")
    http.ListenAndServe(":8080", nil)

    // Wait for goroutine to finish
    wg.Wait()
    fmt.Println("All goroutines finished")
}

