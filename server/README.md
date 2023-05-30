# server

This is a WebSockets go webserver.


## Events (READ)

- `NEW_SESSION`: create a new session
- `GET_DICTIONARY`: send dictionary words back to client
- `CLIENT_READY`: client is ready to start the game, updates this client in session
- `CLIENT_SCORE`: score sent up for the client connected
- `GAME_SUMMARY`: send back game summary to client
