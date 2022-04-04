package session

import (
	"net/http"
	"sync"
	"time"
)

// this map stores the users sessions. For larger scale applications, you can use a database or cache for this purpose
var Sessions sync.Map

func interfaceToStruct(object interface{}) Session {
	session, ok := object.(Session)
	if ok {
		return session
	}
	return session
}

// var Sessions = map[string]Session{}

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
	userSessionInterface, exists := Sessions.Load(sessionToken)
	if !exists {
		return false, 0
	}

	userSessionStruct := interfaceToStruct(userSessionInterface)
	// Позже перенесу в GarbageCollector
	if userSessionStruct.IsExpired() {
		Sessions.Delete(sessionToken)
		return false, 0
	}
	return true, userSessionStruct.ID
}

func LogOutPreviousSession(prev_id int) {
	Sessions.Range(func(key, value interface{}) bool {
		userSession := interfaceToStruct(value)
		if userSession.ID == prev_id {
			Sessions.Delete(key)
		}
		return true
	})
}

func Crear(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:   "session_token",
		MaxAge: -1,
	})
}
