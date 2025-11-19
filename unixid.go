package unixid

import (
	. "github.com/cdvelop/tinystring"
	"github.com/cdvelop/tinytime"
)

const sizeBuf = int32(19)

// lockHandler represents a mutex-like locking mechanism for thread safety
// Typically a sync.Mutex or similar implementation is used
type lockHandler interface {
	Lock()
	Unlock()
}

// userSessionNumber is an interface to obtain the current user's session number
// This is primarily used in WebAssembly environments to uniquely identify client sessions
type userSessionNumber interface {
	// userSessionNumber returns a unique identifier for the current user session
	// e.g., "1" or "2" or "34" or "400" etc.
	userSessionNumber() string
}

// defaultEmptySession provides a default implementation of userSessionNumber that returns an empty string
// This is used in non-WebAssembly environments where session numbers are not needed
type defaultEmptySession struct{}

func (defaultEmptySession) userSessionNumber() string {
	return ""
}

// defaultNoOpMutex provides a mutex implementation that doesn't perform any locking
// This is safe to use in WebAssembly environments which are single-threaded
type defaultNoOpMutex struct{}

func (defaultNoOpMutex) Lock()   {}
func (defaultNoOpMutex) Unlock() {}

// UnixID is the main struct for ID generation and handling
// It contains all configuration and state needed for ID generation
type UnixID struct {
	// userNum is the user session identifier (used in WebAssembly environments)
	userNum string

	// lastUnixNano stores the last generated timestamp to detect collisions
	lastUnixNano int64

	// correlativeNumber is incremented when two IDs would otherwise have the same timestamp
	correlativeNumber int64

	// buf is a pre-allocated buffer to minimize allocations during ID generation
	buf []byte

	// Config holds the external dependencies for the UnixID
	*Config
}

// Config holds the configuration and dependencies for a UnixID instance
type Config struct {
	// Session provides user session numbers in WebAssembly environments
	Session userSessionNumber // e.g., userSessionNumber() string = "1","4","4000" etc.

	// TimeProvider provides time utilities including nanosecond timestamps and date formatting
	TimeProvider tinytime.TimeProvider // Provides UnixNano(), UnixSecondsToDate(), and UnixNanoToTime()

	// syncMutex provides thread safety for concurrent ID generation
	syncMutex lockHandler // e.g., sync.Mutex{}
}

// NewUnixID creates a new UnixID handler with appropriate configuration based on the runtime environment.
//
// For WebAssembly environments (client-side):
// - Requires a userSessionNumber handler to be passed as a parameter
// - Creates IDs with format: "[timestamp].[user_number]" (e.g., "1624397134562544800.42")
// - No mutex is used as JavaScript is single-threaded
//
// For non-WebAssembly environments (server-side):
// - Does not require any parameters
// - Creates IDs with format: "[timestamp]" (e.g., "1624397134562544800")
// - Uses a sync.Mutex for thread safety
//
// IMPORTANT: When integrating with other libraries that also use sync.Mutex,
// you can pass an existing mutex as a parameter to avoid potential deadlocks.
// When an external mutex is provided, this library will use a no-op mutex
// internally to prevent deadlocks when GetNewID is called within a context
// that has already acquired the same mutex. This assumes that external
// synchronization is being handled by the caller.
//
// Parameters:
//   - handlerUserSessionNumber: Optional userSessionNumber implementation (required for WebAssembly)
//   - sync.Mutex or *sync.Mutex: Optional mutex to use instead of creating a new one (server-side only)
//
// Returns:
//   - A configured *UnixID instance
//   - An error if the configuration is invalid
//
// Usage examples:
//
//	// Server-side usage:
//	idHandler, err := unixid.NewUnixID()
//
//	// WebAssembly usage:
//	type sessionHandler struct{}
//	func (sessionHandler) userSessionNumber() string { return "42" }
//	idHandler, err := unixid.NewUnixID(&sessionHandler{})
//
//	// Server-side usage with existing mutex to avoid deadlocks:
//	var mu sync.Mutex
//	idHandler, err := unixid.NewUnixID(&mu)
//
//	// With external mutex, when calling within a locked context:
//	var mu sync.Mutex
//	idHandler, err := unixid.NewUnixID(&mu)
//
//	mu.Lock()
//	defer mu.Unlock()
//	// This won't deadlock because NewUnixID uses a no-op mutex internally
//	// when an external mutex is provided
//	id := idHandler.GetNewID()
func NewUnixID(handlerUserSessionNumber ...any) (*UnixID, error) {
	// The actual implementation is in the build-specific files
	// This function declaration allows for a unified API
	// Implementation details are in unixid_front.go and unixid_back.go
	// and are selected at build time based on the target platform
	return createUnixID(handlerUserSessionNumber...)
}

func configCheck(c *Config) (*UnixID, error) {
	if c == nil {
		return nil, Err(D.Required, D.Configuration, D.Options)
	}

	if c.TimeProvider == nil {
		return nil, Err(D.Required, D.Time, D.Handler)
	}

	// Aseguramos que Session no sea nil (debería estar configurado en createUnixID)
	if c.Session == nil {
		return nil, Err(D.Required, D.Session, D.Handler)
	}

	// Aseguramos que syncMutex no sea nil (debería estar configurado en createUnixID)
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
	currentUnixNano := id.TimeProvider.UnixNano()

	if currentUnixNano == id.lastUnixNano {
		//mientras sean iguales sumar numero correlativo
		id.correlativeNumber++
	} else {
		id.correlativeNumber = 0
	}
	// actualizo la variable unix nano
	id.lastUnixNano = currentUnixNano

	currentUnixNano += id.correlativeNumber

	return Convert(currentUnixNano).String()
}

// GetNewID generates a new unique ID based on Unix nanosecond timestamp.
// In WebAssembly environments, this appends a user session number to the timestamp.
// In server environments, this returns just the Unix nanosecond timestamp.
// Returns a string representation of the unique ID.
func (id *UnixID) GetNewID() string {
	// Aplicamos un bloqueo para garantizar la seguridad del hilo
	id.syncMutex.Lock()
	defer id.syncMutex.Unlock()

	outID := id.unixIdNano()

	// Obtenemos o actualizamos el número de usuario si es necesario
	if id.userNum == "" {
		id.userNum = id.Session.userSessionNumber()
	}

	// Solo añadimos el número de sesión si es válido
	if id.userNum != "" {
		outID += "."
		outID += id.userNum
	}

	return outID
}
