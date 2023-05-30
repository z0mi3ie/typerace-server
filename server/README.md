# server

This is a WebSockets go webserver.


## Connection

- `localhost:8070/` connect to server

## Functions

- on client connect, send generated session id back to client
  - every client will get their own session ID when they connect
  - when all clients in a session disconnect, delete the server session (delete empty sessions)
  - server.CreateSession() string
  - server.DeleteSession(sess string)


## API Messages (READ)

- `CLIENT_READY`: client is ready to start the game, updates this client in session
- `GET_DICTIONARY`: send dictionary words back to client
- `CLIENT_SCORE`: score sent up for the client connected
- `GAME_SUMMARY`: send back game summary to client
