# Simple UDP Server with Packet Reassembly

This Go program implements a simple UDP server designed to receive and reassemble packets. It uses the DNS protocol library `github.com/miekg/dns` to handle DNS messages.

## Code Breakdown

### 1. Imports
The necessary packages are imported, including networking, synchronization, and the DNS library.

### 2. Constants
- **DNSPort**: The port on which the server listens for incoming UDP packets (port 53).
- **TimeoutPeriod**: The duration after which the server will reset its state if no packets are received.

### 3. Main Function
- Resolves the UDP address and starts listening on it.
- Initializes a map to store received packets and a variable to track the expected sequence number.
- Continuously reads packets from clients in a loop and spawns a goroutine for each packet to handle it concurrently.

### 4. handlePacket Function
- This function processes incoming packets.
- It first checks if the packet is a termination signal (indicated by "close").
- It then parses the packet data, including sequence number, packet index, total packets, and the encoded chunk of data.
- If the sequence number is unexpected, it logs an error.
- The encoded chunk is decoded from hex and stored in a map using the packet index.
- An acknowledgment (ACK) message is sent back to the client for each received packet.
- If all expected packets are received, it reassembles them into a complete DNS message and prints it.
- A timeout mechanism resets the state if no packets are received within the specified period.

### 5. Makefile
The Makefile provides commands to build, run, clean, install dependencies, and display Go version information. It allows building the server binary with specific flags to reduce binary size.

## Usage

### 1. Build the Server
Run `make build` to compile the server into a binary named `server`.

### 2. Run the Server
Execute `make run` to start the server.

### 3. Clean Up
Use `make clean` to remove the compiled binary.

### 4. Install Dependencies
Use `make install` to tidy up Go modules.

## Notes
- Ensure you have Go installed and set up correctly on your machine.
- Running this server may require root privileges since it listens on port 53, which is a privileged port.
- The server expects packets formatted in a specific way (with fields separated by `|`) and encoded in hexadecimal format for the data chunk.
- The timeout mechanism resets the state if no packets are received within the specified period, ensuring that stale data does not accumulate.

This implementation is suitable for a basic UDP packet handling scenario with reassembly logic and acknowledgment feedback but may need enhancements for production use, including better error handling and logging.
