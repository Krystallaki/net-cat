// Package server implements the TCP listener and client lifecycle for TCPChat.
package server
import "net"
// Server manages the TCP listener and connected clients.
type Server struct {
	registry *registry
	listener net.Listener
	port string
}

// TODO (Vasiliki): Start(port string) error — bind listener, accept connections, enforce max 10
func (s *Server) Start(port string) error {
s.port = port
listener,err:= net.Listen("tcp",":"+ port)
if err!=nil {
	return err
}
s.listener=listener
s.registry= &registry{
	clients:make(map[*client.Client]struct{})
}
for  {
	conn,err:=s.listener.Accept()
	if err!=nil {
		log.Println("failed to accept connection:",err)
		continue
	}
	if len(s.registry.clients)>=10 {
		conn.Write([]byte("Chat is full. Try again later.\n"))	
		conn.Close()
		continue
	}
	newClient:=&client.Client{Conn:conn,Name:""}
	s.registry.Add(newClient)
	go handleClient(newClient)
}
return nil
}