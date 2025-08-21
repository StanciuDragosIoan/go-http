package main

/*
	fmt → printing output.
	os → for opening files.
	log → for error logging.
	bytes → utility functions for working with byte slices (IndexByte in this case).
*/
import (
	"fmt"
	"log"
	"net"

	"boot.theprimagen.tv/internal/request"
)

func main() {
	/*
		Instead of opening a file, you open a TCP listener.
		:42069 means: listen on port 42069 on all network interfaces.
		This is like running a Node.js net.createServer.
	*/
	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatal("error", "error", err)
	}

	/*
		Accept() blocks until a client connects (e.g. telnet localhost 42069).
		When a client connects, you get a net.Conn.
		net.Conn implements io.ReadWriteCloser, so it works perfectly with your getLinesChannel (just like a file did).
	*/
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("error", "error", err)
		}

		r, err := request.RequestFromReader(conn)

		if err != nil {
			log.Fatal("error", "error", err)
		}

		fmt.Printf("Request line:\n")
		fmt.Printf("- Method: %s\n", r.RequestLine.Method)
		fmt.Printf("- Target: %s\n", r.RequestLine.RequestTarget)
		fmt.Printf("- Version: %s\n", r.RequestLine.HttpVersion)
	}
}

// go run ./cmd/tcplistener | tee requestline.txt
// curl http://localhost:42069/prime/agen
