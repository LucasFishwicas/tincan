package handlers
// Websocket handler functions

import (
    "fmt"
    "net/http"
    "bufio"
    "log"
    "os"
    "sync"

    "github.com/gorilla/websocket"
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
)

func init() {
    // Initialising sendChan
    sendChan = make(chan string)
    // Initialising receiveChan
    receiveChan = make(chan string)
}

// NOT USED -- Handler function to handle requests
func Handle(w http.ResponseWriter, r *http.Request) {
    // Check if request wants to upgrade to Websocket
    if r.Header.Get("Upgrade") == "websocket" {
        WsHandler(w, r)
    } else {
        fmt.Println("Not a websocket upgrade")
    }
}


// Read message sent from user server-side
func ReadSend(wg *sync.WaitGroup) {
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
func ReadReceive(wg *sync.WaitGroup, conn *websocket.Conn, client string) {
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


func WsHandler(w http.ResponseWriter, r *http.Request) {
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
    go ReadSend(&wg)
    go ReadReceive(&wg, conn, client)

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

