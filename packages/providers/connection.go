package providers

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	upgrade = websocket.Upgrader{
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
	pairUP = MatchMake{isBusy: false}
)

const (
	roomIsReady     = 0
	playerHasLeft   = 1
	playerHasJoined = 2
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
		isready     chan bool
		reconnect   chan bool
		isdestroyed chan string
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
		Code int    `json:"code"`
		Gid  string `json:"gid"`
		Data Data   `json:"data"`
	}
	// Data composes the client message
	Data struct {
		Test int `json:"test"`
	}
)

// Init initializes a connection and commits it to memory, afterwards
// the Conn provider will recognize the player using cp.connected[string]
// and the connection will have been appened to the cp.connection stack.
func (cp *ConnProvider) Init(g *gin.Context, created chan bool) error {

	pid := g.Param("id")
	if cp.isConnected(pid) {
		cookie, err := g.Cookie("gid")

		if err != nil {
			return fmt.Errorf("Something went wrong")
		}

		_, ok := Conn.rooms[cookie]

		if ok {
			for v := range Conn.rooms[cookie].player {
				if Conn.rooms[cookie].player[v].pid == pid {
					Conn.rooms[cookie].player[v].client, err = upgrade.Upgrade(g.Writer, g.Request, nil)
					Conn.rooms[cookie].player[v].reconnect <- true
					return nil
				}
			}
		}

		err = fmt.Errorf("User already connected")
		fmt.Println(err)
		return err
	}

	ws, err := upgrade.Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	conn := &connection{
		client:      ws,
		pid:         pid,
		gamePos:     -1,
		reconnect:   make(chan bool),
		isready:     make(chan bool),
		isdestroyed: make(chan string),
	}

	cp.setConn(pid, conn)

	created <- true

	cp.append(conn.pid)

	go pairUP.run(conn)
	return nil
}

// PumpOut pumps out reading and writing from and to our client
func (cp *ConnProvider) PumpOut(pid string, created bool) {
	r := time.Now()

	go writePump(cp.fetch(pid), 0, r)
	go readPump(cp.fetch(pid), 0, r)
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
		return
	}
	m.isBusy = true
	go m.doWork()
	return
}

// this function is responsible for actually processing the matchmaking
func (m *MatchMake) doWork() {
	fmt.Println("Matching")
	for Conn.Size() >= 2 {
		v := make(map[int]*connection)
		c1 := Conn.pop()
		c2 := Conn.pop()

		v[0] = Conn.connected[c1]
		v[0].gamePos = 0

		v[1] = Conn.connected[c2]
		v[1].gamePos = 1

		fmt.Println("Matched!")
		Conn.createRoom(v)
	}

	if Conn.Size() == 0 {
		m.isBusy = false
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
func (cp *ConnProvider) pop() string {
	cp.mu.Lock()
	defer cp.mu.Unlock()
	var conn string

	conn, cp.connections = cp.connections[0], cp.connections[1:]

	return conn
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
		p.isready <- true
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

func closeDoor(gid string) {
	r := Conn.rooms[gid]
	if len(r.player) == 0 {
		delete(Conn.rooms, gid)
		fmt.Printf("ROOM %s is closed!\n", gid)
		r.isdestroyed <- true
	}
}

func serveRoom(r *rooms) {
	for {
		select {
		case <-r.turn:
			r.switchPlayerTurn()

		case index := <-r.playerleft:
			r.player[index].isdestroyed <- r.gid
			go Conn.clearConn(r.player[index].pid, r.gid, index)

		case <-time.After(60 * time.Second):
			r.switchPlayerTurn()

		case <-r.isdestroyed:
			fmt.Println("\n**** Showtime is over! ****")
			return
		}
	}
}

func readPump(c *connection, counter int, lastResponse time.Time) {

	r := &ClientMessage{}
	for {
		err := c.client.ReadJSON(r)
		if err != nil {
			c.client.Close()

			if c.gamePos < 0 {
				c.isdestroyed <- ""
				Conn.clearConn(c.pid, "", 0)
				return
			}

			fmt.Println("Player left, waiting for reconnect")
			select {
			case <-time.After(60 * time.Second):
				Conn.rooms[c.gid].playerleft <- c.gamePos
				return
			case <-c.reconnect:
				fmt.Println("player reconnected")
			}
		}
		if r.Code == 2 {
			if Conn.rooms[c.gid].isItYourTurn(c.pid) {
				Conn.rooms[c.gid].turn <- r.Code
			}
		}
	}
}

func writePump(c *connection, counter int, lastResponse time.Time) {
	fmt.Println("Write pumping...")
	for {
		select {
		case <-c.isready:
			msg := &ClientMessage{
				Code: 0,
				Gid:  c.gid,
				Data: Data{Test: 1},
			}
			fmt.Printf("Hey %s, it's your turn dude", c.pid)
			c.client.WriteJSON(msg)

		case <-c.isdestroyed:
			return

		case <-time.After(50 * time.Second):
			c.client.WriteJSON("ping")
		}
	}
}
