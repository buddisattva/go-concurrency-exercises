//////////////////////////////////////////////////////////////////////
//
// Your video processing service has a freemium model. Everyone has 10
// sec of free processing time on your service. After that, the
// service will kill your process, unless you are a paid premium user.
//
// Beginner Level: 10s max per request
// Advanced Level: 10s max per user (accumulated)
//

package main

import (
	"time"
)

// User defines the UserModel. Use this to check whether a User is a
// Premium user or not
type User struct {
	ID        int
	IsPremium bool
	TimeUsed  int64 // in seconds
}

var timestamps = make(map[*User]chan time.Time)

// HandleRequest runs the processes requested by users. Returns false
// if process had to be killed
func HandleRequest(process func(), u *User) bool {
	if timestamps[u] == nil {
		timestamps[u] = make(chan time.Time, 1)
	}

	accumulateTimeUsed(u)
	process()
	accumulateTimeUsed(u)

	if !u.IsPremium && u.TimeUsed >= 10 {
		return false
	}
	return true
}

func accumulateTimeUsed(u *User) {
	select {
	case t := <-timestamps[u]:
		u.TimeUsed += int64(time.Since(t).Seconds())
	default:
		u.TimeUsed = 0
	}
	timestamps[u] <- time.Now()
	//fmt.Printf("UserID: %d || %d sec elapsed\n", u.ID, u.TimeUsed)
}

func main() {
	RunMockServer()
}
