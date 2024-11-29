package main

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
)

const (
	DNSPort    = ":53"           // DNS server port
	serverAddr = "127.0.0.1"     // DNS server address
	Timeout    = 3 * time.Second // Timeout for receiving ACK
	MaxRetries = 10              // Max retries for each packet
)

func main() {
	var domainName string
	fmt.Print("Enter domain name to query: ")
	fmt.Scanln(&domainName)

	// Build the DNS query for the domain name
	message, err := buildDNSQuery(domainName)
	if err != nil {
		fmt.Println("Error building DNS query:", err)
		return
	}

	addr, err := net.ResolveUDPAddr("udp", serverAddr+DNSPort)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Error dialing UDP:", err)
		return
	}
	defer conn.Close()

	seqNum := rand.Intn(1000)
	fmt.Printf("Session sequence number: %d\n", seqNum)

	totalPackets := (len(message) + 11) / 12
	fmt.Printf("Splitting DNS query into %d packets\n", totalPackets)

	for packetIndex := 0; packetIndex < totalPackets; packetIndex++ {
		start := packetIndex * 12
		end := start + 12
		if end > len(message) {
			end = len(message)
		}
		chunk := message[start:end]

		encodedChunk := hex.EncodeToString(chunk)
		packet := fmt.Sprintf("%d|%d|%d|%s", seqNum, packetIndex, totalPackets, encodedChunk)

		retries := 0
		ackReceived := false

		for retries < MaxRetries && !ackReceived {
			_, err := conn.Write([]byte(packet))
			if err != nil {
				fmt.Println("Error sending packet:", err)
				return
			}
			fmt.Printf("Sent packet %d of %d (sequence number: %d)\n", packetIndex+1, totalPackets, seqNum)

			conn.SetReadDeadline(time.Now().Add(Timeout))
			buffer := make([]byte, 1024)
			n, _, err := conn.ReadFromUDP(buffer)
			if err != nil {
				if strings.Contains(err.Error(), "i/o timeout") {
					retries++
					fmt.Printf("Timeout: No ACK received, resending... (%d/%d)\n", retries, MaxRetries)
				} else {
					fmt.Println("Error receiving ACK:", err)
					return
				}
			} else {
				ackResponse := string(buffer[:n])
				ackParts := strings.Split(ackResponse, "|")
				if len(ackParts) >= 2 && ackParts[0] == fmt.Sprintf("%d", seqNum) && ackParts[1] == fmt.Sprintf("%d", packetIndex) {
					ackReceived = true
					fmt.Printf("ACK received for packet %d\n", packetIndex+1)
				} else {
					fmt.Println("Incorrect ACK received, resending...")
				}
			}
		}

		if !ackReceived {
			fmt.Printf("Failed to receive ACK for packet %d. Aborting.\n", packetIndex+1)
			return
		}
	}

	_, err = conn.Write([]byte(fmt.Sprintf("%d|close|0|0", seqNum)))
	if err != nil {
		fmt.Println("Error sending close message:", err)
		return
	}
	fmt.Println(hex.EncodeToString(message))
	fmt.Println("Session complete. Closing connection.")
}

func buildDNSQuery(domain string) ([]byte, error) {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(domain), dns.TypeA) // Type A query (IPv4 address)
	msg.RecursionDesired = true                  // Set recursion desired flag
	return msg.Pack()                            // Pack the DNS query into byte format
}
