# Smash-Arena
A new take on the now defunct #1 turn-based, multiplayer, browser game ever, Naruto-Arena.

## Tasks:
- [ ] finish UI
- [ ] Implement connection to our game-server
- [ ] Model our data archtecture   
- [ ] Define game data models 

## Data Models

### Character 
**char**

| Key | Value | Description
| --- | ---- | :--- |
| ID | **int** | A unique ID for each character
| Name | **Array[String]** | Character's name
| Profile | **Array[String]** | Flavor text; short introduction

### Start game (SEND)
**startGameReq**

| Key | Value | Description
| --- | ---- | :--- |
| UserID | **String** | A unique ID for each player so our server can identify them
| TeamID | **Array[String]** | Unique identifier for each character the player has on his team

## API endpoints

> GET /character

| parameter | required | info
| --- | --- | --- |
| q | no | this get request will return a JSON with all characters within the game

**Response**

| Key | Value | Description
| --- | ---- | :--- |
| roster | [char[]](#character) | an array containing every character within the game
