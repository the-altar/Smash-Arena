# Smash-Arena
A new take on the now defunct #1 turn-based, multiplayer, browser game ever, Naruto-Arena.

## Tasks:
- [ ] finish UI
- [ ] Implement connection to our game-server
- [ ] Model our data archtecture   
- [ ] Define game data models 

## Data Model - CLIENT
Still under construction
### Start game (SEND)


| Key | Value | Description
| --- | ---- | :--- |
| userID | **String** | A unique ID for each player so our server can identify them
| teamID | **Array[String]** | An array with the ID of each character the player wants to start a game with. The server will use this ID to create a game room
