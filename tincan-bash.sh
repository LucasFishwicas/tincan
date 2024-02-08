#!/bin/bash

echo "Container: Started docker container"

ipAddr="localhost:8080"
function readMessage() {
    read -p "message: " message
    curl -v -G "$ipAddr/http/send" \
        --data-urlencode "message=$message" \
        --data-urlencode "user=$name"
}

curl -v $ipAddr

read -p "http or websocket? " httpWebsocket
if [ $httpWebsocket = "http" ]; then
    read -p "send or receive? " sendReceive
    if [ $sendReceive = "send" ]; then
        read -p "name: " name
        message="enter"
        curl -v -G "$ipAddr/http/send" \
            --data-urlencode "message=$message" \
            --data-urlencode "user=$name"
        while [[ $message != "exit" ]]; do
            readMessage
        done
    elif [ $sendReceive = "receive" ]; then
        watch -n 5 curl -v "$ipAddr/http/receive"
    else
        echo "unrecognised input"
    fi
elif [ $httpWebsocket = "websocket" ]; then
    wscat -c "ws://$ipAddr/ws" \
        --header "Connection: Upgrade" \
        --header "Upgrade: websocket" \
        --header "Host: $ipAddr" \
        --header "Sec-WebSocker-Key: SGVsbG8sIHdvcmxk" \
        --header "Sec-WebSocket-Version: 13"
else
    echo "unrecognised input"
fi

echo "Container: Exiting docker container"
