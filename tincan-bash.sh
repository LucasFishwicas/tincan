#!/bin/bash

echo "Container: Started docker container"
curl -v "localhost:8080"
read -p "send or receive? " response
if [ $response = "send" ]; then
    read -p "name: " name
    curl -v -G "localhost:8080/http/send" --data-urlencode "message=enter" --data-urlencode "user=$name"
    read -p "message: " message
    echo "name: $name, message: $message"
    curl -v -G "localhost:8080/http/send" --data-urlencode "message=$message" --data-urlencode "user=$name"
    while [[ $message != "exit" ]]; do
        read -p "message: " message
        echo "name: $name, message: $message"
        curl -v -G "localhost:8080/http/send" --data-urlencode "message=$message" --data-urlencode "user=$name"
    done
elif [ $response = "receive" ]; then
    watch -n 2 curl -v "localhost:8080/http/receive"
else
    echo "unrecognised input"
fi
echo "Container: Exiting docker container"
# wscat -c ws://localhost:8080/ws --header "Connection: Upgrade" --header "Upgrade: websocket" --header "Host: localhost:8080" --header "Sec-WebSocker-Key: SGVsbG8sIHdvcmxk" --header "Sec-WebSocket-Version: 13"
