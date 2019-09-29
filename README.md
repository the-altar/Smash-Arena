# Smash-Arena
A new take on the now defunct #1 turn-based, multiplayer, browser game ever, Naruto-Arena.

## Tasks:
- [ ] finish UI
- [ ] Implement connection to our game-server
- [ ] Model our data archtecture   
- [ ] Define game data models 

# Data Models

## **Character** 
### charClient
This is the data sent back to the client when they make a request for a specific character

| Key | Value | Description |
| --- | ---- | :--- |
| ID | **int** | A unique ID for each character
| Name | **Array[String]** | Character's name
| Profile | **Array[String]** | Flavor text; short introduction

## charServer 

| Key | Value | Description |
| --- | ---- | :--- |
| ID | **int** | A unique identifier
| Name | **string** | Character's name
| Health | **int** | Character's current health
| Skills | [**Map{int}Skill**](#Skills) | Flavor text; short introduction

## **Skills**
### Skill
| Key | Value | Description | Client-side
| --- | ---- | :--- | :---: |
| ID | **int** | unique identifier | yes
| Name | **bool** | skill's name; also unique | yes
| Desc | **int** | Description of what the skill does | Yes
| Effect | [**map{String}Effect**](#Skill-Effects)| No

## **Skill Effects**
Every struct defined below are lumped together by an *Effect* interface 
### Damage 

| Key | Value | Description | Client-side
| --- | ---- | :--- | :---: |
| Value | **int** | How much damage is dealt | yes
| Tick | **bool** | If false, damage will be dealt every turn, otherwise it'll be dealt only after effect's duration ends | yes
| Duration | **int** | How long this effect will last. | yes



## Game
### **startGameReq**
This is the data the server expects when the client wants to start a new game

| Key | Value | Description
| --- | ---- | :--- |
| UserID | **String** | A unique ID for each player so our server can identify them
| TeamID | **Array[String]** | Unique identifier for each character the player has on his team

## API endpoints

### get character
> GET    /character

| parameter | required | info
| --- | --- | --- |
| q | no | this get request will return a JSON with all characters within the game

**Response**

| Key | Value | Description
| --- | ---- | :--- |
| roster | [char[]](#charClient) | an array containing client-relevant information about every character within the game

### newgame
> POST   /newgame

| key | Value | required
| --- | --- | --- |
| playerID | **string** | yes
| teamID |**[]string** | yes

**Response** 

| key | Value | Description
| --- | --- | --- |
| gid | **string** | A unique identifier for a game room

### arena
> GET    /arena

| parameter | required | info
| --- | --- | --- |
| q | no | This request will open a websocket between the client and the server

**Response**

| key | Value | Description
| --- | --- | --- |
| code | **int** | 1 if sucessful, 0 otherwise


