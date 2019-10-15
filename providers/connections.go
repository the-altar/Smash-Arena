package providers

import (
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
)

var (
	upgrade = websocket.Upgrader{}
	// Conn manages our connections
	Conn = &ConnProvider{
		connections: make([]*connection, 0),
		connected:   make(map[string]interface{}),
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
		client  *websocket.Conn
		pid     string
		isready chan string
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
		connected   map[string]interface{}
		rooms       map[string]*rooms
	}
)

// Init initializes a connection and commits it to memory, afterwards
// the Conn provider will recognize the player using cp.connected[string]
// and the connection will have been appened to the cp.connection stack.
func (cp *ConnProvider) Init(c echo.Context) error {
	pid := c.Param("id")
	if cp.isConnected(pid) {
		return echo.ErrForbidden
	}

	ws, err := upgrade.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}

	cp.setConn(pid, true)

	conn := &connection{ws, pid, make(chan string)}

	cp.append(conn)

	go pairUP.Run(conn)

	// clients waiting for the goroutine to finish wait here
	cp.setConn(pid, <-conn.isready)

	return nil
}

// Run matchmaking function
func (m *MatchMake) Run(conn *connection) {
	if m.isBusy {
		return
	}
	m.isBusy = true
	go m.doWork()
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

// this function is responsible for actually processing the matchmaking
func (m *MatchMake) doWork() {
	for Conn.Size() >= 2 {
		c1 := Conn.pop()
		c2 := Conn.pop()

		v := [2]*connection{c1, c2}

		Conn.createRoom(v)
	}

	if Conn.Size() == 0 {
		m.isBusy = false
		return
	}

	Conn.keepalive()
	time.Sleep(10 * time.Second)
	go m.doWork()
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

// checks if clients are still connected. If a connection goes missing, remove all traces of it from the server
func (cp *ConnProvider) keepalive() {
	for v := range cp.connections {

		go func(c *connection) {
			err := c.client.WriteMessage(websocket.PingMessage, []byte("ping"))
			if err != nil {
				log.Fatal(err)
			}
		}(cp.connections[v])

	}
}

// set the connection status. Bool and String values are acccepted in the second parameter.
func (cp *ConnProvider) setConn(id string, val interface{}) {
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
