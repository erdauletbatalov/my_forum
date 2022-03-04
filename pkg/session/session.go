package session

import (
	"net/http"
	"time"

	"github.com/erdauletbatalov/forum.git/pkg/models"
)

// this map stores the users sessions. For larger scale applications, you can use a database or cache for this purpose
var Sessions = map[string]Session{}

// each session contains the username of the user and the time at which it expires
type Session struct {
	Email  string
	Expiry time.Time
}

// we'll use this method later to determine if the session has expired
func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}

func IsSession(r *http.Request) (bool, *models.User) {
	user := &models.User{}
	c, err := r.Cookie("session_token")
	if err != nil {
		return false, user
	}
	sessionToken := c.Value
	userSession, exists := Sessions[sessionToken]
	if !exists {
		return false, user
	}
	user.Email = userSession.Email
	if userSession.IsExpired() {
		delete(Sessions, sessionToken)
		return false, user
	}
	return true, user
}
