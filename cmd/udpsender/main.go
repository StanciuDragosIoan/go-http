// This is a simple UDP client written in Go.
// It connects to a UDP server running on localhost:42069,
// then continuously waits for user input from the terminal,
// and sends whatever the user types to the server.

package main

import (
	"bufio" // for reading user input from the console
	"fmt"   // for printing messages to the console
	"net"   // for UDP networking
	"os"    // for access to standard input (keyboard)
)

func main() {
	// --- Step 1: Define the server address ---
	// Resolve the string "127.0.0.1:42069" into a usable UDP address structure.
	// 127.0.0.1 = localhost, 42069 = port number
	remoteAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:42069")
	if err != nil {
		fmt.Println("Error resolving remote address:", err)
		return // stop the program if we cannot resolve the address
	}

	// --- Step 2: Connect to the server ---
	// net.DialUDP creates a UDP connection object.
	// The first parameter (local address) is nil, so the OS picks a free local port automatically.
	// The second parameter (remoteAddr) is the server we want to send messages to.
	conn, err := net.DialUDP("udp", nil, remoteAddr)
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
		return // stop the program if connection fails
	}
	// Ensure the connection is closed automatically when main() ends
	defer conn.Close()

	// --- Step 3: Setup a reader for user input ---
	// Wrap standard input (keyboard) in a buffered reader so we can easily read lines of text.
	reader := bufio.NewReader(os.Stdin)

	// Print a message so the user knows the client is ready.
	fmt.Println("UDP client ready. Type messages:")

	// --- Step 4: Infinite loop for reading input and sending ---
	for {
		// Print a prompt character to indicate we're waiting for input
		fmt.Print("> ")

		// Read a line of text from the user (until Enter is pressed)
		// ReadString includes the newline character '\n' at the end.
		line, err := reader.ReadString('\n')
		if err != nil {
			// If reading from stdin fails, show the error but keep going
			fmt.Println("Error reading input:", err)
			continue
		}

		// Send the typed line to the UDP server
		_, err = conn.Write([]byte(line))
		if err != nil {
			// If writing fails, show the error but keep going
			fmt.Println("Error writing to UDP:", err)
			continue
		}
	}
}
