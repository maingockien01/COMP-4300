# File

- `Assignment 1 - Report.pdf` Report/Screenshots of my assignment
- `COMP4300-A1.postman_collection.json` Postman collection for testing Rest API
- `Dockerfile.server` Dockerfile for creating images of server. 
- Folder `cmd` contain entry points
- Folder `bin` (after run make) contain executable files.
- Folder `pkf`:
        - `client` chat client
        - `models` program domain models
        - `services`
                - `chat/server` chat server
                - `protocol` websocket protocols
        - `websocket` package for websocket
# Run

## Chat Server
```
$ make buildChatService
$ HOST=localhost PORT=8010 SERVER_NAME=chat1 bin/chatService
```

### With docker
If you have docker and want to run the server with Docker. Recommend to run normally without docker for the sake of simplicity.
```
$ docker build -f Dockerfile.server -t "discovery_service:Dockerfile.server" .
$ docker run --publish 8010:8080 --network=host  server:Dockerfile.server 
```

## Chat client
```
$ make buildChatClient
$ CHAT_URL=127.0.0.1:8010 USER_ID=user_client_1 USER_NAME=Kevin_client_1 bin/chatClient
```

## Notes
- I set the port for server is 8010 but you feel free to change it to what you want.
- The host is localhost if you run on your own computer. However, you need to change it to correct host if you run on aviary.

# Chat server
## API

### Room

`GET /rooms` to get the list of all rooms

`POST /room` to create a new room
Body: 
```
{
        "name":"[RoomName]"
}
```


### User

`GET /users/room/{RoomName}` to get list of users in [RoomName] room

`POST /users/room/{RoomName}` to add users into room
Body:
```
{
        "id": "[UserId]"
        "name": "[UserName]"
}
```

`DELETE /users/room/{RoomName}` to rremove users from room
Body:
```
{
        "id": "[UserId]"
        "name": "[UserName]"
}
```

`POST /user` to register users with server (client automatically register on start up)
Body:
```
{
        "id": "[UserId]"
        "name": "[UserName]"
}
```

### Message
`GET /messages/{roomName}` to get all messages in the room

`POST /messages/{roomName}` to add message into room chat
Body:
```
{
        {"sender":"[UserId]", "content": "[Message]", "room": "[RoomName]"}
}
```

## WebSocket
- connection: WebSocket
- path: `/chat`

1. Gretting protocol
Used to register new user with server or register websocket for that user so client will receive message from server. Currently, server only support one socket for a user. If users use 2 client, olny 1 will receive message.

```
{"metadata": {"version": 1, "from": "[UserId]", "direction": "greeting", "type": "user"},"data": {"id": "[UserId]", "name": "[UserName"}}
```

2. Message protocol (client receive)
```
{"metadata":{"version":1,"from":"[ServerName]","direction":"push","type":"message"},"data":[{"sender":"[UserId]","content":"[Content]","room":"[RoomName]", "position": "[Position]", "timeAt":"[Time]"}]}
```

3. Message protocol (client update)
```
{"metadata":{"version":1,"from":"[UserId]","direction":"update","type":"message"},"data":[{"sender":"[UserId]","content":"[Content]","room":"[RoomName]"}]}
```


# Chat client

- type `\help` to see all commands
- the USER_ID is like "email" used to register account with server. Client would register automatically when it starts up. However, you can use postman to register your own user and interact with server. 
- client keep track of the room user is working on. Use `\room_join` to switch or join other room.
- `\message [message]` will send message to the room that user is worknig on currently.
- I used "github.com/gorilla/websocket" for implemnting client websocket. 

# Testing

## Postman
I add the collection of API request I used to test server api. 

## Client
I test client by interacting with it mostly.