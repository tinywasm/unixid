//go:build !wasm

package unixid

import (
	"sync"
)

// createUnixID for server implementation
func createUnixID(params ...any) (*UnixID, error) {
	c := &Config{
		Session:   &defaultEmptySession{},
		syncMutex: &sync.Mutex{},
	}

	externalMutexProvided := false

	for _, param := range params {
		switch p := param.(type) {
		case *sync.Mutex:
			externalMutexProvided = true
		case sync.Mutex:
			externalMutexProvided = true
		case userSessionNumber:
			c.Session = p
		}
	}

	if externalMutexProvided {
		c.syncMutex = &defaultNoOpMutex{}
	}

	return configCheck(c)
}
