package unixid

import (
	. "github.com/tinywasm/fmt"
)

// Parse parses an ID string and extracts its components.
// It first validates the ID format, then extracts the timestamp and optional user number.
//
// Parameters:
//   - id: The ID string to parse (e.g., "1624397134562544800" or "1624397134562544800.42")
//
// Returns:
//   - timestamp: The timestamp portion as int64
//   - userNum: The user number portion as string (empty if not present)
//   - err: An error if the ID format is invalid or parsing fails
func (u *UnixID) Parse(id string) (timestamp int64, userNum string, err error) {
	// Primero valida el formato
	if err := u.Validate(id); err != nil {
		return 0, "", err
	}

	// Encuentra el índice del punto (si existe)
	point_index := len(id)
	for i, char := range id {
		if char == '.' {
			point_index = i
			break
		}
	}

	// Extrae la parte del timestamp
	timestamp_str := id[:point_index]

	// Extrae el user number si existe
	if point_index < len(id) {
		userNum = id[point_index+1:]
	}

	// Convierte el timestamp a int64
	timestamp, er := Convert(timestamp_str).Int64()
	if er != nil {
		return 0, "", Err("format", "invalid")
	}

	return timestamp, userNum, nil
}
