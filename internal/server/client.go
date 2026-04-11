// Package server implements the TCPChat server logic.
package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

const welcomeBanner = "Welcome to TCP-Chat!\n" +
	"         _nnnn_\n" +
	"        dGGGGMMb\n" +
	"       @p~qp~~qMb\n" +
	"       M|@||@) M|\n" +
	"       @,----.JM|\n" +
	"      JS^\\__/  qKL\n" +
	"     dZP        qKRb\n" +
	"    dZP          qKKb\n" +
	"   fZP            SMMb\n" +
	"   HZM            MMMM\n" +
	"   FqM            MMMM\n" +
	" __| \".        |\\dS\"qML\n" +
	" |    `.       | `' \\Zq\n" +
	"_)      \\.___.,|     .'\n" +
	"\\____   )MMMMMP|   .'\n" +
	"     `-'       `--'\n" +
	"[ENTER YOUR NAME]: "

// Client represents a connected chat client.
type Client struct {
	Conn net.Conn
	Name string
}

// sendWelcome sends the Linux banner and name prompt to the client.
func sendWelcome(conn net.Conn) {
	fmt.Fprint(conn, welcomeBanner)
}

// readName reads a non-empty name from the client connection.
// If the client sends an empty name, it re-prompts until a valid name is given.
func readName(conn net.Conn) (string, error) {
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		name := strings.TrimSpace(scanner.Text())
		if name != "" {
			return name, nil
		}
		fmt.Fprint(conn, "[ENTER YOUR NAME]: ")
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("connection closed before name was provided")
}
