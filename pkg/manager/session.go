package manager

// SESSION defines a session struct
type SESSION struct {
	ID       int
	Username string
}

// Sessions maps authenticated users
var Sessions = map[string]SESSION{}

// SetSession saves a session's value
func SetSession(sid string, ID int, username string) {
	Sessions[sid] = SESSION{ID: ID, Username: username}
}

// GetSession retrives a session's value if it exists; also returns false if it doesn't
func GetSession(sid string) (SESSION, bool) {
	val, ok := Sessions[sid]
	if ok {
		return val, true
	}
	return val, false
}

// DestroySession removes session
func DestroySession(sid string) {
	delete(Sessions, sid)
}
