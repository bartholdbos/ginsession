package ginsession

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/gin-gonic/gin"
	"io"
	"net/url"
	"time"
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
	AddSession(ID string) (Session, error) //Create a new session with ID
	GetSession(ID string) (Session, error) //Get session by ID
	DelSession(ID string) error            //Delete session by ID
	ClearSessions(lifetime int64)          //Clear inactive sessions
}

//Create a new Manager
func CreateManager(name string, lifetime int64, providername string) (*Manager, error) {
	provider, ok := providers[providername]
	if ok {
		return &Manager{name: name, lifetime: lifetime, provider: provider}, nil
	} else {
		return nil, errors.New("session.CreateManager: Provider Unknown")
	}
}

func generateID() (string, error) {
	b := make([]byte, 32)

	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func (manager *Manager) SessionInit(c *gin.Context) (session Session, err error) {
	var ID string

	cookie, err1 := c.Cookie(manager.name)
	if err1 != nil || cookie == "" {
		ID, err = generateID()
		if err != nil {
			return
		}

		session, err = manager.provider.AddSession(ID)
		c.SetCookie(manager.name, url.QueryEscape(ID), int(manager.lifetime), "/", "", false, true)
	} else {
		ID, err = url.QueryUnescape(cookie)
		if err != nil {
			return
		}

		session, err = manager.provider.GetSession(ID)
	}

	return
}

func (manager *Manager) SessionDestroy(c *gin.Context) (err error) {
	var ID string

	cookie, err1 := c.Cookie(manager.name)
	if err1 != nil || cookie == "" {
		return
	} else {
		ID, err = url.QueryUnescape(cookie)
		if err != nil {
			return
		}

		manager.provider.DelSession(ID)
		c.SetCookie(manager.name, "", -1, "/", "", false, true)
	}

	return
}

func (manager *Manager) SessionClear() {
	manager.provider.ClearSessions(manager.lifetime)
	time.AfterFunc(time.Duration(manager.lifetime), func() { manager.SessionClear() })
}

func Register(name string, provider Provider) error {
	if provider != nil {
		return errors.New("session.Register: Provider is nil")
	}

	if _, dup := providers[name]; dup {
		return errors.New("session.Register: Provider already registered")
	}

	providers[name] = provider

	return nil
}
