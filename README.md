 ---  IN DEVELOPMENT  ---

A bi-directional chat app designed to be used through the terminal

-----------------------------------------

## Run the front end

1. Within the frontEnd directory (containing Dockerfile and bash script), build the Docker image with the command:
`docker build -t tincan .`

2. Run a Docker container based off the newly created image, tincan:
`docker run -it --rm tincan`

3. The container should immediately run the bash script which will attempt to connect to the server, followed by a prompt for your preferred protocol (http or websocket)

### websocket
- Responding to the prompt with "websocket" will call a wscat command to connect your terminal to the server and enable instant 2 way communication with the server

### http
- Responding to the prompt with "http" will result in a seconday prompt for "send" or "receive":
- "receive" will call a watch command which calls curl on the /http/receive endpoint of the server, repeatedly pulling messages at regularly intervals
- "send" will prompt for a message and then call a curl command on /http/send enpoint, submitting this message to the server (to be pulled at /http/receive)

- To effectively communicate through the http protocol, it is recommended to open 2 separate terminal windows, 1 for "receive" and the other for "send"
