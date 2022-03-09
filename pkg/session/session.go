package session

import (
	"net/http"
	"time"
)

// this map stores the users sessions. For larger scale applications, you can use a database or cache for this purpose
var Sessions = map[string]Session{}

// each session contains the username of the user and the time at which it expires
type Session struct {
	ID     int
	Expiry time.Time
}

// we'll use this method later to determine if the session has expired
func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

// IsSession checks if user and server has simmilar session. It returns true and user's credentials if
// so and false and empty user if not
func IsSession(r *http.Request) (bool, int) {
	c, err := r.Cookie("session_token")
	if err != nil {
		return false, 0
	}
	sessionToken := c.Value
	userSession, exists := Sessions[sessionToken]
	if !exists {
		return false, 0
	}

	// Позже перенесу в GarbageCollector
	if userSession.IsExpired() {
		delete(Sessions, sessionToken)
		return false, 0
	}
	return true, userSession.ID
}
