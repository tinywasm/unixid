package unixid

// SetNewID generates a new unique ID and assigns it to target.
//
// In WebAssembly environments, IDs include a user session number as a suffix (e.g., "1624397134562544800.42").
// In server environments, IDs are just the timestamp (e.g., "1624397134562544800").
//
// This function is thread-safe in server-side environments.
//
// Example:
//
//	// Setting a string variable
//	var id string
//	idHandler.SetNewID(&id)
//
//	// Setting a struct field
//	type User struct{ ID string }
//	user := User{}
//	idHandler.SetNewID(&user.ID)
func (id *UnixID) SetNewID(target *string) {
	*target = id.NewID()
}
