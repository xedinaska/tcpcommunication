package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const (
	cmdSTOP = "STOP"
	cmdSEND = "SEND"
)

func main() {
	serverHostPtr := flag.String("server-host", "localhost", "use -server-host to provide server host for client app (localhost by default)")
	serverPortPtr := flag.Int("server-port", 3333, "use -server-port to provide server port for client app (3333 by default)")
	flag.Parse()

	server := fmt.Sprintf("%s:%d", *serverHostPtr, *serverPortPtr)

	conn, err := net.Dial("tcp", server)
	if err != nil {
		log.Fatalf("[FATAL]host=%s,port=%d; failed to connect to provided server: `%s`", *serverHostPtr, *serverPortPtr, err.Error())
	}

	defer conn.Close()

	log.Printf("[INFO]client app successfully connected to server %s", server)

	//listen for server messages in separate routine
	go func() {
		buff := make([]byte, 1024)
		for {
			n, err := conn.Read(buff)
			if err != nil {
				log.Printf("[ERROR]failed to read server message: `%s`", err.Error())
				continue
			}

			msg := string(buff[:n])
			log.Printf("[DEBUG]received message from server: `%s`", msg)

			if msg == cmdSTOP {
				log.Printf("[INFO]STOP signal received from server; stopping client..")
				os.Exit(0)
				break
			}
		}
	}()

	//reading commands from stdin
	for msg := range readStdin(os.Stdin) {
		//parse received message. Should contains 2 parts: COMMAND:MESSAGE
		input := strings.Split(msg, ":")
		if len(input) != 2 {
			log.Printf("[ERROR]invalid input syntax, expected `COMMAND:message_text_here`")
			continue
		}

		cmd, message := input[0], input[1]

		//check parsed COMMAND. If STOP - exit client. If SEND - send message to server. If smth else: show error
		switch cmd {
		case cmdSTOP:
			log.Printf("[INFO]STOP signal received; stopping client..")
			os.Exit(0)
		case cmdSEND:
			conn.Write([]byte(message))
			log.Printf("[DEBUG]send message `%s` to `%s`", message, server)
		default:
			log.Printf("[ERROR]invalid command received, available input: `STOP:`, `SEND:MESSAGE_BODY_HERE`")
		}
	}
}

//readStdin used to run the goroutine that reading stdin messages
func readStdin(r io.Reader) <-chan string {
	lines := make(chan string)
	go func() {
		defer close(lines)
		scan := bufio.NewScanner(r)
		for scan.Scan() {
			lines <- scan.Text()
		}
	}()
	return lines
}
