package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
)

const (
	DNSPort       = ":53"
	TimeoutPeriod = 120 * time.Second
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", DNSPort)
	if err != nil {
		fmt.Println("Error resolving address:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Error starting UDP server:", err)
		return
	}
	defer conn.Close()

	fmt.Println("Server listening on", DNSPort)

	packetMap := make(map[int]string)
	var expectedSeqNum string
	var wg sync.WaitGroup

	for {
		buffer := make([]byte, 1024)
		n, clientAddr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error reading from UDP:", err)
			continue
		}

		wg.Add(1)
		go handlePacket(conn, buffer[:n], clientAddr, &wg, &packetMap, &expectedSeqNum)
	}
	wg.Wait()
}

func handlePacket(conn *net.UDPConn, packetData []byte, clientAddr *net.UDPAddr, wg *sync.WaitGroup, packetMap *map[int]string, expectedSeqNum *string) {
	defer wg.Done()

	timeoutTimer := time.NewTimer(TimeoutPeriod)

	packetDataStr := string(packetData)
	fmt.Printf("\nReceived packet: %s\n", packetDataStr)

	parts := strings.Split(packetDataStr, "|")
	if len(parts) < 4 {
		fmt.Println("Malformed packet, skipping...")
		return
	}

	seqNum := parts[0]
	packetIndex := parts[1]
	totalPacketsStr := parts[2]
	encodedChunk := parts[3]

	if packetIndex == "close" {
		fmt.Println("Received termination signal. Closing session.")
		*packetMap = make(map[int]string)
		*expectedSeqNum = ""
		timeoutTimer.Stop()
		wg.Wait()
		return
	}

	packetIndexInt, err := strconv.Atoi(packetIndex)
	if err != nil {
		fmt.Println("Error parsing packet index:", err)
		return
	}

	totalPackets, err := strconv.Atoi(totalPacketsStr)
	if err != nil {
		fmt.Println("Error parsing total packets:", err)
		return
	}

	if *expectedSeqNum == "" {
		*expectedSeqNum = seqNum
	} else if seqNum != *expectedSeqNum {
		fmt.Printf("Unexpected sequence number: got %s, expected %s.\n", seqNum, *expectedSeqNum)
		return
	}

	decodedChunk, err := hex.DecodeString(encodedChunk)
	if err != nil {
		fmt.Println("Error decoding packet content:", err)
		return
	}
	(*packetMap)[packetIndexInt] = string(decodedChunk)
	fmt.Printf("Received and stored packet %d: %s\n", packetIndexInt, string(decodedChunk))

	ackMsg := fmt.Sprintf("%s|%d", seqNum, packetIndexInt)
	_, err = conn.WriteToUDP([]byte(ackMsg), clientAddr)
	if err != nil {
		fmt.Println("Error sending ACK:", err)
		return
	}
	fmt.Printf("Sent ACK for packet %d\n", packetIndexInt)

	if len(*packetMap) == totalPackets {
		fmt.Println("\nReceived packets in order:")
		var indices []int
		for idx := range *packetMap {
			indices = append(indices, idx)
		}
		sort.Ints(indices)

		var reassembledData []byte
		for _, idx := range indices {
			reassembledData = append(reassembledData, []byte((*packetMap)[idx])...)
		}

		var msg dns.Msg
		err := msg.Unpack(reassembledData)
		if err != nil {
			fmt.Println("Malformed DNS packet:", err)
		} else {
			fmt.Println("\nReassembled DNS packet:", msg.String())
		}

		*packetMap = make(map[int]string)
		*expectedSeqNum = ""
		timeoutTimer.Stop()
	}

	timeoutTimer.Reset(TimeoutPeriod)

	select {
	case <-timeoutTimer.C:
		*packetMap = make(map[int]string)
		*expectedSeqNum = ""
		return
	}
}
