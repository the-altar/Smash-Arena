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
		connections: make([]string, 0),
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
		isdestroyed chan string
	}

	rooms struct {
		sync   sync.Mutex
		player []*connection
		gid    string
	}
	// ConnProvider manages our connections
	ConnProvider struct {
		mu          sync.Mutex
		connections []string
		connected   map[string]*connection
		rooms       map[string]*rooms
	}
)

// Init initializes a connection and commits it to memory, afterwards
// the Conn provider will recognize the player using cp.connected[string]
// and the connection will have been appened to the cp.connection stack.
func (cp *ConnProvider) Init(g *gin.Context, created chan bool) error {

	pid := g.Param("id")
	if cp.isConnected(pid) {
		return g.Error(nil)
	}

	ws, err := upgrade.Upgrade(g.Writer, g.Request, nil)
	if err != nil {
		fmt.Println(err)
		return err
	}

	conn := &connection{ws, pid, "", make(chan string), make(chan string)}

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
	for Conn.Size() >= 2 {
		v := make([]*connection, 0)
		c1 := Conn.pop()
		c2 := Conn.pop()

		v = append(v, Conn.connected[c1], Conn.connected[c2])

		Conn.createRoom(v)
		fmt.Println(Conn.Size())
	}

	if Conn.Size() == 0 {
		m.isBusy = false
		return
	}
	time.Sleep(20 * time.Second)
	go m.doWork()
	return
}

// creates a room with 2 players each mapped by a UUID string
func (cp *ConnProvider) createRoom(r []*connection) {
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
func (cp *ConnProvider) clearConn(pid string, gid string) {
	cp.mu.Lock()
	defer cp.mu.Unlock()

	_, ok := cp.rooms[gid]

	if ok {
		for v := range Conn.rooms[gid].player {
			if Conn.rooms[gid].player[v].pid == pid {
				Conn.rooms[gid].player = append(Conn.rooms[gid].player[:v], Conn.rooms[gid].player[v+1:]...)
				break
			}
		}
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

	fmt.Println("CURRENT STATUS: ")
	fmt.Println("CONNECTIONS:", len(Conn.connected))
	fmt.Println("ROOMS: ", len(Conn.rooms))
	fmt.Println("QUEUE: ", len(Conn.connections))
	fmt.Printf("-------------\n\n")
	return
}

func closeDoor(gid string) {
	Conn.rooms[gid].sync.Lock()
	defer Conn.rooms[gid].sync.Unlock()

	r := Conn.rooms[gid]
	if len(r.player) == 0 {
		delete(Conn.rooms, gid)
	}

}

func readPump(c *connection, counter int, lastResponse time.Time) {
	type res struct {
		data string
	}
	r := &res{}
	for {
		err := c.client.ReadJSON(r)
		if err != nil {
			c.isdestroyed <- c.gid
			return
		}
	}
}

func writePump(c *connection, counter int, lastResponse time.Time) {
	fmt.Println("Write pumping...")
	for {
		select {

		case c.gid = <-c.isready:
			fmt.Println("CONNECTED")

		case <-c.isdestroyed:
			fmt.Println("Im gonna need you gone")
			err := c.client.Close()
			if err != nil {
				fmt.Println("Err: ", err)
			}
			Conn.clearConn(c.pid, c.gid)

			return

		case <-time.After(54 * time.Second):
			err := c.client.WriteJSON("ping")
			if err != nil {
				fmt.Println("CLIENT HAS LEFT")
			}
		}
	}
}
