package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	var globalHeader, packetHeader []byte

	file, err := os.Open("net.cap")
	if err != nil {
		log.Fatal(err)
	}

	// Read in pcap global header
	globalHeader = make([]byte, 24)
	globalCount, err := file.Read(globalHeader)
	if err != nil {
		log.Fatal(err)
	}

	if globalCount != 24 {
		log.Fatal("pcap global header not correct length")
	}

	fmt.Printf("pcap global header: %x\n", globalHeader)

	// Handle packets
	var packetNum int
	for {
		packetHeader = make([]byte, 16)
		count, err := file.Read(packetHeader)
		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		if count != 16 {
			log.Fatal("packet header not correct length")
		}

		fmt.Printf("packet %d header: %x\n", packetNum, packetHeader)

		packetLength := binary.LittleEndian.Uint16(packetHeader[8:11])
		untruncatedPacketLength := binary.LittleEndian.Uint16(packetHeader[12:15])
		fmt.Printf("%d %d\n", packetLength, untruncatedPacketLength)

		if packetLength != untruncatedPacketLength {
			log.Fatal("packet is truncated")
		}

		packetData := make([]byte, packetLength)
		count, err = file.Read(packetData)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Packet %d Data: %x\n", packetNum, packetData)
		packetNum++
	}
}
