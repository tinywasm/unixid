//go:build wasm

package unixid

// createUnixID for WASM implementation
func createUnixID(handlerUserSessionNumber ...any) (*UnixID, error) {
	c := &Config{
		Session:   &defaultEmptySession{},
		syncMutex: &defaultNoOpMutex{},
	}

	for _, u := range handlerUserSessionNumber {
		if usNumber, ok := u.(userSessionNumber); ok {
			c.Session = usNumber
		}
	}

	return configCheck(c)
}
