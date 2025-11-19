package unixid

import . "github.com/cdvelop/tinystring"

// UnixNanoToStringDate converts a Unix nanosecond timestamp ID to a human-readable date string.
// This function accepts a string representation of a Unix nanosecond timestamp (with or without a user number suffix)
// and returns a formatted date string in the format "2006-01-02 15:04" (year-month-day hour:minute).
//
// Parameters:
//   - unixNanoStr: String representation of a Unix timestamp in nanoseconds (e.g. "1624397134562544800" or "1624397134562544800.42")
//
// Returns:
//   - A formatted date string (e.g. "2021-06-22 15:32")
//   - An error if the input is invalid or the conversion fails
//
// Example:
//
//   dateStr, err := handler.UnixNanoToStringDate("1624397134562544800")
//   if err != nil {
//     // handle error
//   }
//   fmt.Println(dateStr) // e.g. "2021-06-22 15:32"
func (u *UnixID) UnixNanoToStringDate(unixNanoStr string) (string, error) {

	// Convert string ID to int64 (validating it in the process)
	unixNano, err := ValidateID(unixNanoStr)
	if err != nil {
		return "", err
	}

	if u.TimeProvider == nil {
		return "", Err("TimeProvider nil")
	}

	// Format the nanosecond timestamp to "YYYY-MM-DD HH:MM:SS"
	dateTime := u.TimeProvider.FormatDateTime(unixNano)
	// Return only date and time without seconds (HH:MM format)
	if len(dateTime) >= 16 {
		return dateTime[:16], nil // "2006-01-02 15:04"
	}
	return dateTime, nil
}

// https://chat.openai.com/c/4af98def-f8d9-4095-bf31-deaaad84c094
