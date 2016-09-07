package ginsession

import ()

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

//Create a new Manager
func CreateManager(name string, lifetime int64, provider Provider) (*Manager, error) {
	return nil, nil
}
