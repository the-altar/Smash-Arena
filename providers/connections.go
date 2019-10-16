package providers

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	upgrade = websocket.Upgrader{}
	// Conn manages our connections
	Conn = &ConnProvider{
		connections: make([]*connection, 0),
		connected:   make(map[string]*connection),
		rooms:       make(map[string]*rooms),
	}
	// PairUP connections that are available
	pairUP = MatchMake{isBusy: false}
)

type (
	// MatchMake is responsible for matching clients together once there are enough connections
	MatchMake struct {
		isBusy bool
	}

	connection struct {
		client      *websocket.Conn
		pid         string
		gid         string
		isready     chan string
		isdestroyed chan bool
	}

	rooms struct {
		sync   sync.Mutex
		player [2]*connection
		gid    string
	}
	// ConnProvider manages our connections
	ConnProvider struct {
		lock        sync.Mutex
		connections []*connection
		connected   map[string]*connection
		rooms       map[string]*rooms
	}
)

// Init initializes a connection and commits it to memory, afterwards
// the Conn provider will recognize the player using cp.connected[string]
// and the connection will have been appened to the cp.connection stack.
func (cp *ConnProvider) Init(g *gin.Context, created chan bool) error {
	fmt.Println("-> INIT HAS BEGUN")
	pid := g.Param("id")
	if cp.isConnected(pid) {
		fmt.Println("User already connected")
		return g.Error(nil)
	}

	ws, err := upgrade.Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		fmt.Println("Failed to upgrade connection")
		return err
	}

	conn := &connection{ws, pid, "", make(chan string), make(chan bool)}

	cp.setConn(pid, conn)

	created <- true

	cp.append(conn)

	go pairUP.run(conn)
	fmt.Println("-> CLOSED INIT")
	return nil
}

// PumpOut pumps out reading and writing from and to our client
func (cp *ConnProvider) PumpOut(pid string, created bool) {

	fmt.Println("Pumping out for ", pid)
	go writePump(cp.fetch(pid), 0)
	fmt.Println("locked")
	return
}

// Size function returns how many connections we have
func (cp *ConnProvider) Size() int {
	cp.lock.Lock()
	size := len(cp.connections)
	cp.lock.Unlock()

	return size
}

// checks if a user is Connected to our server or not. Returns true or false
func (cp *ConnProvider) isConnected(pid string) bool {
	cp.lock.Lock()
	defer cp.lock.Unlock()
	_, ok := cp.connected[pid]
	return ok
}

// Run matchmaking function
func (m *MatchMake) run(conn *connection) {
	fmt.Println("-> MM.RUN HAS BEGUN")
	if m.isBusy {
		fmt.Println("-> MM.RUN HAS BEEN CLOSED")
		return
	}
	m.isBusy = true
	go m.doWork()
	fmt.Println("MM.RUN HAS BEEN CLOSED")
	return
}

// this function is responsible for actually processing the matchmaking
func (m *MatchMake) doWork() {
	fmt.Println("-> doWork BEGUN")
	for Conn.Size() >= 2 {
		c1 := Conn.pop()
		c2 := Conn.pop()

		v := [2]*connection{c1, c2}

		Conn.createRoom(v)
	}

	if Conn.Size() == 0 {
		m.isBusy = false
		fmt.Println("-> doWork HAS BEEN CLOSED")
		return
	}
	fmt.Println("Attempted to matchmake")
	time.Sleep(10 * time.Second)
	go m.doWork()
	return
}

// creates a room with 2 players each mapped by a UUID string
func (cp *ConnProvider) createRoom(r [2]*connection) {
	id, err := uuid.NewUUID()
	if err != nil {
		panic("couldn't create ID")
	}
	gid := id.String()
	rs := &rooms{player: r, gid: gid}
	cp.rooms[gid] = rs

	rs.player[0].isready <- gid
	rs.player[1].isready <- gid

	return
}

// set the connection status. Bool and String values are acccepted in the second parameter.
func (cp *ConnProvider) setConn(id string, val *connection) {
	cp.lock.Lock()
	cp.connected[id] = val
	cp.lock.Unlock()
}

// appends a connection to the connection stack
func (cp *ConnProvider) append(conn *connection) {
	cp.lock.Lock()
	cp.connections = append(cp.connections, conn)
	cp.lock.Unlock()
}

func (cp *ConnProvider) fetch(pid string) *connection {
	c := cp.connected[pid]
	return c
}

// Pops a connection from the connection stack
func (cp *ConnProvider) pop() *connection {
	cp.lock.Lock()
	defer cp.lock.Unlock()
	var conn *connection

	conn, cp.connections = cp.connections[0], cp.connections[1:]

	return conn
}

// Removes a connection from stack with a given "pid"
func (cp *ConnProvider) remove(pid string) {
	cp.lock.Lock()
	defer cp.lock.Unlock()

	for v := range cp.connections {
		if cp.connections[v].pid == pid {
			cp.connections = append(cp.connections[:v], cp.connections[v+1:]...)
			break
		}
	}

	delete(cp.connected, pid)
}

func writePump(c *connection, counter int) {
	fmt.Println("Write pumping...")
	for {
		select {
		case c.gid = <-c.isready:
			fmt.Println(c.pid, " in room ", c.gid, " is ready to be closed")
			c.client.Close()
			return
		case <-time.After(54 * time.Second):
			fmt.Println("Timedout")
			c.client.WriteJSON("ping")
		}
	}
}
