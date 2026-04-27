// Package server implements the TCP listener and client lifecycle for TCPChat.
// It accepts incoming connections, enforces a maximum client limit,
// and dispatches each connection to its own goroutine for concurrent handling.
package server

import (
	"log"
	"net"
	"net-cat/internal/client"
	"net-cat/internal/messaging"
)

// Server manages the TCP listener and connected clients.
// It holds the active listener, the port it is bound to,
// and a thread-safe registry of all currently connected clients.
type Server struct {
	registry *registry
	listener net.Listener
	port     string
}

// Start binds the TCP listener on the given port, initialises the client registry,
// and accepts incoming connections in a loop. Each connection is registered and
// handled in its own goroutine. Returns an error if the listener cannot be created.
func (s *Server) Start(port string) error {
	s.port = port
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}
	s.listener = listener
	s.registry = &registry{
		clients: make(map[*client.Client]struct{}),
	}
	history := &messaging.History{}
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("failed to accept connection:", err)
			continue
		}
		if len(s.registry.All()) >= 10 {
			if _, err := conn.Write([]byte("Chat is full. Try again later.\n")); err != nil {
				log.Println("failed to notify full chat:", err)
			}
			conn.Close()
			continue
		}

		newClient := &client.Client{Conn: conn, Name: ""}
		s.registry.Add(newClient)
		go s.handleClient(newClient, history)
	}
	return nil
}
func (s *Server) handleClient(newClient *client.Client, history *messaging.History) {
	defer s.registry.Remove(newClient)
	defer newClient.Conn.Close()
	client.SendWelcome(newClient.Conn)
	name, err := client.ReadName(newClient.Conn)
	if err != nil {
		return
	}
	newClient.Name = name
	history.Replay(newClient.Conn)
	client.NotifyJoin(s.registry.All(), newClient)
	messaging.RunMessageLoop(newClient, s.registry.All, history)
	client.NotifyLeave(s.registry.All(), newClient)
}
