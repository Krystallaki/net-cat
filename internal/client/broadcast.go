package client

// Broadcast sends msg to all clients in the provided slice except exclude.
// exclude may be nil for server announcements such as join/leave.
func Broadcast(clients []*Client, msg string, exclude *Client) {
	// TODO (Krysta)
}
