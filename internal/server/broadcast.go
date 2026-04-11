package server

// TODO (Krysta): broadcast sends msg to all registered clients except exclude.
// exclude may be nil (e.g. for server announcements like join/leave).
func broadcast(msg string, exclude *Client) {
}
