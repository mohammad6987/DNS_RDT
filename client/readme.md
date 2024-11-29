# DNS RDT Client

This Go program implements a simple UDP client that sends DNS queries to a specified server. The client splits the DNS query into multiple packets, sends them to the server, and waits for acknowledgments (ACK) for each packet.

## Features

- Sends DNS queries over UDP.
- Splits large queries into smaller packets.
- Implements retry logic for packet transmission.
- Receives and verifies ACKs from the server.
- Closes the session gracefully after completing the query.

## Requirements

- Go 1.15 or later
- github.com/miekg/dns package (automatically handled by Go modules)

## Usage

### Building the Client

To build the client, use the following command:

```bash
make build
```

This will compile the Go code into a binary named client.

### Running the Client

To run the client after building, use:

```bash
make run
```

The client will prompt you to enter a domain name to query. After entering the domain name, it will send a DNS query to the server.

### Cleaning Up

To remove the compiled binary and temporary files, run:

```bash
make clean
```

### Installing Dependencies

To install or tidy up Go module dependencies, run:

```bash
make install
```

### Displaying Information

To view the Go version and environment details, use:

```bash
make info
```

### Help

For a list of available Makefile targets and their descriptions, run:

```bash
make help
```

## Code Overview

### Main Functionality

1. **Input**: The program prompts the user to input a domain name.
2. **DNS Query Construction**: It constructs a DNS query using the buildDNSQuery function.
3. **UDP Connection**: Establishes a UDP connection to the server at 127.0.0.1:53.
4. **Packet Sending**: Splits the DNS query into packets of 12 bytes each and sends them to the server.
5. **ACK Handling**: Waits for ACKs from the server for each packet and implements a retry mechanism in case of timeouts.
6. **Session Closure**: Sends a close message to indicate the end of the session.

### BuildDNSQuery Function

This function constructs a DNS query for a given domain name using the miekg/dns library and returns it in byte format.

## Example

When you run the client, it will look like this:
```bash
Enter domain name to query: example.com
Session sequence number: 123
Splitting DNS query into 2 packets
Sent packet 1 of 2 (sequence number: 123)
ACK received for packet 1
Sent packet 2 of 2 (sequence number: 123)
ACK received for packet 2
Session complete. Closing connection.
```
## Notes

- Make sure you have a DNS server running and accessible at 127.0.0.1:53 before running this client.
