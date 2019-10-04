# Game State Message

```json
{
    client: string,
    code: int,
    gameState: {
        opponent: string,
        foes: {
            char: {
                id: int,
                health: int,
                skills: {}
            }
        },

        friends: {
            char:{
                id: int,
                health: int,
                skills: {
                    id: int
                }
            }
        }
    }

}
```
