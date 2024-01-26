package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func startClient(port string) {
	time.Sleep(time.Second * 2)
	conn, err := net.Dial("tcp", port)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println("Connection established")

	conn.Write([]byte("C"))

	buf := make([]byte, 2048)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print("> ")
		scanner.Scan()
		cmd := scanner.Text()

		if cmd == "EXIT" {
			break
		}

		conn.Write([]byte(cmd))

		n, err := conn.Read(buf)
		if err != nil {
			log.Fatal(err)
			return
		}

		res := buf[:n]
		fmt.Println(string(res))
	}

	conn.Close()
}
