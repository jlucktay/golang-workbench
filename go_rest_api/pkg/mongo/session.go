// Package mongo interfaces with MongoDb for us.
package mongo

import (
	"gopkg.in/mgo.v2"
)

// Session holds a MongoDb session.
type Session struct {
	session *mgo.Session
}

// Open connects to a local MongoDb server.
func (s *Session) Open() error {
	var err error
	s.session, err = mgo.Dial("127.0.0.1:27017")
	if err != nil {
		return err
	}
	s.session.SetMode(mgo.Monotonic, true)
	return nil
}

// Copy returns a pointer to the session.
func (s *Session) Copy() *mgo.Session {
	return s.session.Copy()
}

// Close will close a non-nil session.
func (s *Session) Close() {
	if s.session != nil {
		s.session.Close()
	}
}
