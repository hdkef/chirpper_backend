## Chirpper
Chirpper is twitter clone app which highly depends on websocket technology to deliver realtime communication. It uses angular framework as frontend and golang as backend.
This is the backend. Basically contains a lot of goroutines that handles websocket connection.

# Auth method 
Because of websocket connection can't have bearer header by default, the token is being sent via payload and if the token is invalid or expired, the websocket get disconnected. When the user's refresh the page, then the autologin will sent the token to server and server respons with CLEAR BEARER response header so that the client is forced to log out.

# TODO
Follow method haven't been implemented because Firebase currently doesn't support array manipulation.
Until then, this feature will not be created.