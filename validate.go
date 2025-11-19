package unixid

import (
	. "github.com/cdvelop/tinystring"
)

// Validate validates the format of an ID string without parsing it.
// It handles IDs in both server format (just timestamp) and client format (timestamp.userNumber).
//
// Parameters:
//   - id: The ID string to validate (e.g., "1624397134562544800" or "1624397134562544800.42")
//
// Returns:
//   - error: nil if valid, error describing the issue if invalid
//
// Validation rules:
//   - The ID must not be empty
//   - The ID must contain only digits and at most one decimal point
//   - The ID must not start or end with a decimal point
//   - The timestamp portion (before the decimal point) must be valid
func (u *UnixID) Validate(id string) error {
	msg_invalid := Err(D.Character, D.Invalid, D.Not, D.Supported)

	if len(id) == 0 {
		return msg_invalid
	}

	// No debe comenzar ni terminar con punto
	if id[0] == '.' || id[len(id)-1] == '.' {
		return msg_invalid
	}

	var point_count int
	for _, char := range id {
		if char == '.' {
			point_count++
			if point_count > 1 {
				return Err(D.Format, D.Invalid, D.Found, D.More, D.Point)
			}
		} else if char < '0' || char > '9' {
			return msg_invalid
		}
	}

	return nil
}
