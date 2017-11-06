package webserver

// Error is our default error layout when returning to the client
type Error struct {
	Err interface{} `json:"err"`
}
