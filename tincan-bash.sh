#!/bin/bash

echo "Container: Started docker container"
curl "localhost:8080"
read -p "send or receive? " response
if [ $response = "send" ]; then
    echo "send!"
    read -p "name: " name
    curl -G "localhost:8080/http/send" --data-urlencode "message=enter" --data-urlencode "user=$name"
    read -p "message: " message
    echo "name: $name, message: $message"
    curl -G "localhost:8080/http/send" --data-urlencode "message=$message" --data-urlencode "user=$name"
    while [[ $message != "exit" ]]; do
        read -p "message: " message
        echo "name: $name, message: $message"
        curl -G "localhost:8080/http/send" --data-urlencode "message=$message" --data-urlencode "user=$name"
    done
elif [ $response = "receive" ]; then
    echo "receive"
    watch -n 2 curl "localhost:8080/http/receive"
else
    echo "unrecognised input"
fi
echo "Container: Exiting docker container"
# wscat -c ws://localhost:8080/ws --header "Connection: Upgrade" --header "Upgrade: websocket" --header "Host: localhost:8080" --header "Sec-WebSocker-Key: SGVsbG8sIHdvcmxk" --header "Sec-WebSocket-Version: 13"
