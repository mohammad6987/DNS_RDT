# UDP Packet Transmission: Client-Server Application

This project demonstrates a simple client-server application that sends and receives packets over UDP, simulating a DNS query transmission. The client splits a DNS query into multiple packets, sends them to the server, and waits for an acknowledgment (ACK) for each packet. The server receives the packets, sends ACKs back, and processes them when all packets are received. The connection is terminated when the client sends a termination signal.

## Features

- **Client**: 
  - Accepts a domain name, builds a DNS query, splits it into multiple packets.
  - Sends packets to the server and waits for ACKs for each packet.
  - Terminates the session by sending a "close" signal after all packets are acknowledged.

- **Server**:
  - Receives packets, processes them, and sends an ACK for each packet.
  - Once all packets are received, the server outputs the complete message and resets for the next session.
  - If the "close" signal is received from the client, it terminates the current session.

## Files

- **client.go**: The client program that builds the DNS query, splits it into packets, sends them, and waits for ACKs.
- **server.go**: The server program that receives packets, acknowledges each, and handles termination.
- **Makefile**: Automation for building, running, and cleaning the project.
  
## Prerequisites

Ensure you have the following installed:

- Go (1.18 or higher)
- `make` (for running the Makefile)

## Setup

1. **Clone the repository**:

   ```bash
   git clone https://github.com/your-username/udp-client-server.git
   cd udp-client-server
