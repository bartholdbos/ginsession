package ginsession

import (
	"errors"
)

var providers = make(map[string]Provider)

type Manager struct {
	name     string
	provider Provider
	lifetime int64
}

type Session interface {
	Set(key string, value interface{}) error //Set session value for key
	Get(key string) (interface{}, error)     //Return session value for key
	Del(key string) error                    //Delete session value for key
	ID() string                              //Return session ID
}

type Provider interface {
	AddSession(ID string) (Session, err) //Create a new session with ID
	GetSession(ID string) (Session, err) //Get session by ID
	DelSession(ID string) err            //Delete session by ID
	ClearSessions(lifetime int64)        //Clear inactive sessions
}

//Create a new Manager
func CreateManager(name string, lifetime int64, providername string) (*Manager, error) {
	provider, ok := providers[providername]
	if ok {
		return &Manager{name: name, lifetime: lifetime, provider: providername}, nil
	} else {
		return nil, errors.New("Unknown Provider")
	}
}
