package arena

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/the-altar/Smash-Arena/pkg/game"
)

var (
	// Upgrade Upgrades a get request to a websocket
	Upgrade = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	// Conn manages our connections
	Conn = &ConnProvider{
		connections: make([]string, 0),
		connected:   make(map[string]*connection),
		rooms:       make(map[string]*rooms),
	}
	// PairUP connections that are available
	PairUP = MatchMake{isBusy: false}
)

type (
	// MatchMake is responsible for matching clients together once there are enough connections
	MatchMake struct {
		isBusy bool
	}

	connection struct {
		client      *websocket.Conn
		gamePos     int
		pid         string
		gid         string
		isready     chan int
		reconnect   chan bool
		isdestroyed chan string
		Arena       *game.Arena
	}

	rooms struct {
		sync        sync.Mutex
		player      map[int]*connection
		turnCount   int
		gid         string
		isdestroyed chan bool
		playerleft  chan int
		turn        chan int
	}
	// ConnProvider manages our connections
	ConnProvider struct {
		mu          sync.Mutex
		connections []string
		connected   map[string]*connection
		rooms       map[string]*rooms
	}

	// ClientMessage is the data structure sent to the client and what we expect to receive
	ClientMessage struct {
		Code     int        `json:"code"`
		Gid      string     `json:"gid"`
		GameData game.Arena `json:"gameData"`
		TeamID   []int      `json:"teamId"`
	}
)

// Begin initializes a connection and commits it to memory, afterwards
// the Conn provider will recognize the player using cp.connected[string]
// and the connection will have been appened to the cp.connection stack.
func (cp *ConnProvider) Begin(g *gin.Context, created chan bool) error {

	fmt.Println("Attempted to connect")
	pid := g.Param("id")
	if cp.isConnected(pid) {
		cookie, err := g.Cookie("gid")

		fmt.Println(cookie)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Already connected but gid doesn't exist")
			return fmt.Errorf("Connected but gid doesn't exist")
		}

		_, ok := Conn.rooms[cookie]

		if ok {
			for v := range Conn.rooms[cookie].player {
				if Conn.rooms[cookie].player[v].pid == pid {
					Conn.rooms[cookie].player[v].client, err = Upgrade.Upgrade(g.Writer, g.Request, nil)
					Conn.rooms[cookie].player[v].reconnect <- true
					return nil
				}
			}
		}

		err = fmt.Errorf("User already connected")
		fmt.Println(err)
		return err
	}

	ws, err := Upgrade.Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	conn := &connection{
		client:      ws,
		pid:         pid,
		gamePos:     -1,
		reconnect:   make(chan bool),
		isready:     make(chan int),
		isdestroyed: make(chan string),
		Arena:       nil,
	}

	cp.append(conn.pid)
	created <- true
	cp.setConn(pid, conn)
	go PairUP.run(conn)
	return nil
}

// PumpOut pumps out reading and writing from and to our client
func (cp *ConnProvider) PumpOut(pid string, created bool) {
	go writePump(cp.fetch(pid))
	go readPump(cp.fetch(pid))
	return
}

// Size function returns how many connections we have
func (cp *ConnProvider) Size() int {
	cp.mu.Lock()
	size := len(cp.connections)
	cp.mu.Unlock()

	return size
}

// checks if a user is Connected to our server or not. Returns true or false
func (cp *ConnProvider) isConnected(pid string) bool {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	_, ok := cp.connected[pid]
	return ok
}

// Run matchmaking function
func (m *MatchMake) run(conn *connection) {
	if m.isBusy {
		fmt.Println("Worker is busy! Please be patient")
		return
	}
	m.isBusy = true
	go m.doWork()
	return
}

// this function is responsible for actually processing the matchmaking
func (m *MatchMake) doWork() {
	fmt.Println("Worker: Yo, Im working.")
	for Conn.Size() >= 2 {
		fmt.Println("Matching")
		v := make(map[int]*connection)
		c1 := Conn.pop(0)
		c2 := Conn.pop(0)

		v[0] = Conn.connected[c1]
		v[0].gamePos = 0

		v[1] = Conn.connected[c2]
		v[1].gamePos = 1

		v[0].Arena.Enemies = v[1].Arena.Allies
		v[1].Arena.Enemies = v[0].Arena.Allies

		fmt.Println("Matched!")
		Conn.createRoom(v)
	}

	if Conn.Size() == 0 {
		m.isBusy = false
		fmt.Println("Killed matching process")
		return
	}
	time.Sleep(5 * time.Second)
	go m.doWork()
	return
}

// creates a room with 2 players each mapped by a UUID string
func (cp *ConnProvider) createRoom(r map[int]*connection) {
	id, err := uuid.NewUUID()

	if err != nil {
		panic("couldn't create ID")
	}

	gid := id.String()

	for conn := range r {
		r[conn].gid = gid
	}

	rs := &rooms{
		player:      r,
		gid:         gid,
		turnCount:   0,
		turn:        make(chan int),
		isdestroyed: make(chan bool),
		playerleft:  make(chan int),
	}

	cp.rooms[gid] = rs

	go serveRoom(rs)
	rs.turn <- 0
	return
}

// set the connection status. Bool and String values are acccepted in the second parameter.
func (cp *ConnProvider) setConn(id string, val *connection) {
	cp.mu.Lock()
	cp.connected[id] = val
	cp.mu.Unlock()
}

// appends a connection to the connection stack
func (cp *ConnProvider) append(conn string) {
	cp.mu.Lock()
	cp.connections = append(cp.connections, conn)
	cp.mu.Unlock()
}

func (cp *ConnProvider) fetch(pid string) *connection {
	c := cp.connected[pid]
	return c
}

// Pops a connection from the connection stack
func (cp *ConnProvider) pop(counter int) string {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	conn := cp.connections[counter]
	if cp.connected[conn].Arena != nil {
		cp.connections = cp.connections[counter+1:]
		return conn
	}

	return cp.pop(counter + 1)
}

// removes player from server
func (cp *ConnProvider) clearConn(pid string, gid string, pos int) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	_, ok := cp.rooms[gid]
	if ok {
		delete(cp.rooms[gid].player, pos)
		closeDoor(gid)
	} else {
		for v := range cp.connections {
			if cp.connections[v] == pid {
				cp.connections = append(cp.connections[:v], cp.connections[v+1:]...)
				break
			}
		}
	}

	delete(cp.connected, pid)

	fmt.Printf("-------------\n\n")
	fmt.Println("CURRENT STATUS: ")
	fmt.Println("CONNECTIONS:", len(Conn.connected))
	fmt.Println("ROOMS: ", len(Conn.rooms))
	fmt.Println("QUEUE: ", len(Conn.connections))
	fmt.Printf("-------------\n\n")
	return
}

func (r *rooms) countUp() {
	r.sync.Lock()
	r.turnCount++
	r.sync.Unlock()
}

func (r *rooms) switchPlayerTurn() {
	p, ok := r.player[r.playerTurn()]
	if ok {
		p.isready <- 1
	} else {
		fmt.Println("player went missing")
	}

	r.countUp()
}

func (r *rooms) playerTurn() int {
	r.sync.Lock()
	defer r.sync.Unlock()

	counter := r.turnCount % 2
	return counter
}

func (r *rooms) isItYourTurn(pid string) bool {
	r.sync.Lock()
	defer r.sync.Unlock()
	if r.player[r.turnCount%2].pid != pid {
		return true
	}
	return false
}
