package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
    //"os/signal"
    "sync"

	"github.com/gorilla/websocket"
)




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
    

    // Start the server on port 8080
    fmt.Println("Server listening on port 8080")
    http.ListenAndServe(":8080", nil)

    // Wait for goroutine to finish
    wg.Wait()
    fmt.Println("All goroutines finished")
}

