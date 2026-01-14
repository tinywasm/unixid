package unixid

import (
	. "github.com/tinywasm/fmt"
	tinytime "github.com/tinywasm/time"
)

const sizeBuf = int32(19)

// lockHandler represents a mutex-like locking mechanism for thread safety
type lockHandler interface {
	Lock()
	Unlock()
}

// userSessionNumber is an interface to obtain the current user's session number
type userSessionNumber interface {
	userSessionNumber() string
}

// defaultEmptySession provides a default implementation of userSessionNumber
type defaultEmptySession struct{}

func (defaultEmptySession) userSessionNumber() string {
	return ""
}

// defaultNoOpMutex provides a mutex implementation that doesn't perform any locking
type defaultNoOpMutex struct{}

func (defaultNoOpMutex) Lock()   {}
func (defaultNoOpMutex) Unlock() {}

// UnixID is the main struct for ID generation and handling
type UnixID struct {
	userNum           string
	lastUnixNano      int64
	correlativeNumber int64
	buf               []byte
	*Config
}

// Config holds the configuration for a UnixID instance
type Config struct {
	Session   userSessionNumber
	syncMutex lockHandler
}

// NewUnixID creates a new UnixID handler with appropriate configuration based on the runtime environment.
func NewUnixID(handlerUserSessionNumber ...any) (*UnixID, error) {
	return createUnixID(handlerUserSessionNumber...)
}

func configCheck(c *Config) (*UnixID, error) {
	if c == nil {
		return nil, Err(D.Required, D.Configuration, D.Options)
	}

	if c.Session == nil {
		return nil, Err(D.Required, D.Session, D.Handler)
	}

	if c.syncMutex == nil {
		return nil, Err(D.Required, D.Sync, "Mutex")
	}

	return &UnixID{
		userNum:           "",
		lastUnixNano:      0,
		correlativeNumber: 0,
		buf:               make([]byte, 0, sizeBuf),
		Config:            c,
	}, nil
}

func (id *UnixID) unixIdNano() string {
	currentUnixNano := tinytime.Now()

	if currentUnixNano == id.lastUnixNano {
		id.correlativeNumber++
	} else {
		id.correlativeNumber = 0
	}
	id.lastUnixNano = currentUnixNano
	currentUnixNano += id.correlativeNumber

	return Convert(currentUnixNano).String()
}

// GetNewID generates a new unique ID based on Unix nanosecond timestamp (UTC).
func (id *UnixID) GetNewID() string {
	id.syncMutex.Lock()
	defer id.syncMutex.Unlock()

	outID := id.unixIdNano()

	if id.userNum == "" {
		id.userNum = id.Session.userSessionNumber()
	}

	if id.userNum != "" {
		outID += "."
		outID += id.userNum
	}

	return outID
}
