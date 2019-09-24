# Smash-Arena
A new take on the now defunct #1 turn-based, multiplayer, browser game ever, Naruto-Arena.

## Server side
### Packages (go modules)
### Engine
### Server 

## Client side
There's a handful of things yet to be done, client-side, such as:
- [ ] finish UI
- [ ] Implement connection to our game-server
- [ ] Model our data archtecture   

## Data Model
Still under construction
### Start game

| Key  |  Value | Description
| --- | ---- |
| userID | **String** | A unique ID for each player so our server can identify them
| team  | **Array[String]** | An array with the ID of each character the player wants to start a game with. The server will use this ID to create a game room  