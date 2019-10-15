package providers

import (
	"fmt"
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
		connected:   make(map[string]bool),
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
		connected   map[string]bool
		rooms       map[string]*rooms
	}
)

func (m *MatchMake) doWork() {
	fmt.Println("Matchmaking.")
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
	go Conn.keepalive(60 * time.Second)
	time.Sleep(10 * time.Second)
	go m.doWork()
}

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

func (cp *ConnProvider) keepalive(timeout time.Duration) {
	for i := range cp.connections {
		err := cp.connections[i].client.WriteJSON(websocket.PingMessage)
		if err != nil {
			fmt.Println("error: ", err)
			cp.connections[i].client.Close()
		}
	}
}

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

	cp.connected[pid] = true
	conn := &connection{ws, pid, make(chan string)}

	cp.lock.Lock()
	cp.connections = append(cp.connections, conn)
	cp.lock.Unlock()

	go pairUP.Run(conn)
	return nil
}

// Run matchmaking function
func (m *MatchMake) Run(conn *connection) {
	if m.isBusy {
		fmt.Println(<-conn.isready)
		return
	}

	m.isBusy = true
	go m.doWork()
	fmt.Println(<-conn.isready)
}

// Pops a connection from cp.connections[]
func (cp *ConnProvider) pop() *connection {
	cp.lock.Lock()
	defer cp.lock.Unlock()
	var conn *connection

	conn, cp.connections = cp.connections[0], cp.connections[1:]

	return conn
}

// Size function returns how many connections we have
func (cp *ConnProvider) Size() int {
	cp.lock.Lock()
	size := len(cp.connections)
	cp.lock.Unlock()

	return size
}

func (cp *ConnProvider) isConnected(pid string) bool {
	cp.lock.Lock()
	defer cp.lock.Unlock()
	_, ok := cp.connected[pid]
	return ok
}
