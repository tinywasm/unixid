//go:build wasm
// +build wasm

package unixid

import "github.com/tinywasm/time"

// createUnixID para WASM ahora usa time.TimeProvider
func createUnixID(handlerUserSessionNumber ...any) (*UnixID, error) {
	t := time.NewTimeProvider()

	c := &Config{
		Session:      &defaultEmptySession{},
		TimeProvider: t,
		syncMutex:    &defaultNoOpMutex{},
	}

	for _, u := range handlerUserSessionNumber {
		if usNumber, ok := u.(userSessionNumber); ok {
			c.Session = usNumber
		}
	}

	return configCheck(c)
}
