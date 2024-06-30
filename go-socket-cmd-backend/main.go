package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:5555")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go func() {
			read_connection(&conn)
		}()
	}
}

func read_connection(conn *net.Conn) {
	for {
		buffer := make([]byte, 4096)
		_, err := bufio.NewReader(*conn).Read(buffer)
		if err != nil {
			log.Fatal(err)
		}

		// trimming the buffer of rest of it.
		command_buffer := strings.TrimRightFunc(string(buffer),
			func(c rune) bool {
				if !(c > ' ' && c < '~') {
					return true
				} else {
					return false
				}
			})

		result, err := exec.Command("bash", "-c", command_buffer).CombinedOutput()
		if err != nil {
			log.Fatal(err)
		}

		wrote, err := fmt.Fprint(*conn, string(result))
		if err != nil {
			log.Fatal(err)
		} else if wrote < len(result) {
			log.Fatal(err)
		}
	}
}
