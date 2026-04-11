package chat

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

// SendWelcome sends the Linux banner and name prompt to the client.
func SendWelcome(conn net.Conn) {
	fmt.Fprint(conn, welcomeBanner)
}

// ReadName reads a non-empty name from the client connection.
// If the client sends an empty name, it re-prompts until a valid name is given.
func ReadName(conn net.Conn) (string, error) {
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
