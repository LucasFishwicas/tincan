package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/gorilla/websocket"
)




var (
    // Set upgrader variable to access Upgrader method
    upgrader = websocket.Upgrader{
        ReadBufferSize: 1024,
        WriteBufferSize: 1024,
    }

    // Create channel to ensure send before closure
    now_exit = make(chan int, 1)
    // Create channel for sending server-side to client
    sendChan = make(chan string,1)
    // Create channel for receiving client-side to server
    receiveChan = make(chan string,1)
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
        if len(strings.Fields(text)) != 0 {
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
            receiveChan <- "client:// left the chat"
            return
        }

        if len(strings.Fields(string(message))) != 0 {
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
    fmt.Println("client://",client,"entered the chat")
    fmt.Println("")

    // Send a welcome message
    err = conn.WriteMessage(websocket.TextMessage, []byte("tincan:// You're in!"))
    if err != nil {
        log.Println("Error sending message:", err)
        return
    }
    err = conn.WriteMessage(websocket.TextMessage, []byte(" "))
    if err != nil {
        log.Println("Error sending message:", err)
        return
    }

    // Launch goroutines for reading and writing
    wg.Add(2)
    go readSend(&wg)
    go readReceive(&wg, conn)

    // Handle server closure using Ctrl-C
    go exit_session(conn)


    // Eternally loop and check channel messages
    for {
        select {
        case received := <-receiveChan:
            if string(received) == "client:// left the chat" {
                fmt.Println("client://",client," left the chat")
                fmt.Println("")
            } else {
                fmt.Print("-> ",string(received))
            }
        case sent := <-sendChan:
            if sent != "tincan:// Server Shutting Down" {
                sent = "-> "+sent
            }
            if err := conn.WriteMessage(websocket.TextMessage,[]byte(sent)); err != nil {
                log.Println("Error writing message:",err)
                return
            }
            if sent == "tincan:// Server Shutting Down" {
                now_exit <- 1
            }
        
        }
    }
}


// Reserved to handle http requests
func httpHandler() {
    // Handle Regular HTTP Request
    // Sams code can go here...
    fmt.Println("Not a websocket upgrade")
}




// My not-so-graceful attempt at a graceful shutdown
func exit_session(conn *websocket.Conn) {
    exit := make(chan os.Signal, 1)
    signal.Notify(exit, syscall.SIGINT, syscall.SIGQUIT)

    <-exit
    sendChan <- "tincan:// Server Shutting Down"
    <-now_exit
    os.Exit(0)
}






func main() {

    // Define a handler function for incoming requests
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w,"Welcome to tincan\nThe single-line CLI chat service\n") 
    })


    // Handle initiation of Websocket
    http.HandleFunc("/ws",wsHandler)
    

    // Start the server on port 8080
    fmt.Println("Server listening on port 8080")
    fmt.Println("")
    http.ListenAndServe(":8080", nil)


    // Wait for goroutine to finish
    wg.Wait()
    fmt.Println("All goroutines finished")
}

