package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
)

const (
	IPv4EtherType uint16 = 0x0008
	IPv6EtherType uint16 = 0xDD86
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
	var firstEtherType uint16
	httpData := make(map[uint32][]byte)
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

		packetLength := binary.LittleEndian.Uint16(packetHeader[8:11])
		untruncatedPacketLength := binary.LittleEndian.Uint16(packetHeader[12:15])

		if packetLength != untruncatedPacketLength {
			log.Fatal("packet is truncated")
		}

		// Parse Ethernet Headers

		// MAC addresses
		macDestination := make([]byte, 6)
		count, err = file.Read(macDestination)
		if err != nil {
			log.Fatal(err)
		}

		macSource := make([]byte, 6)
		count, err = file.Read(macSource)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("MAC Source/Destination addresses: %x %x\n", macSource, macDestination)

		// EtherType field
		etherTypeData := make([]byte, 2)
		count, err = file.Read(etherTypeData)
		if err != nil {
			log.Fatal(err)
		}

		etherType := binary.LittleEndian.Uint16(etherTypeData)

		// make sure it's IPv4 or IPv6
		if etherType != IPv4EtherType && etherType != IPv6EtherType {
			log.Fatal("The IP datagram must be IPv4 or IPv6")
		}

		if packetNum == 0 {
			firstEtherType = etherType
		}

		if packetNum > 0 && etherType != firstEtherType {
			log.Fatal("Not all IP Datagrams have the same format")
		}

		// IP Packet

		// Internet Header Length (IHL)
		iPHeader := make([]byte, 20)
		count, err = file.Read(iPHeader)
		if err != nil {
			log.Fatal("Error reading IP header")
		}

		// fmt.Printf("Internet Header Length: %x\n", iPHeader[0]&15)

		// Total IP Length
		// totalIPLength := binary.BigEndian.Uint16(iPHeader[2:4])
		// fmt.Printf("IP Length: %d\n", totalIPLength)

		// Source IP Address
		// sourceIPAddr := binary.BigEndian.Uint32(iPHeader[12:16])
		// fmt.Printf("Source IP Address: %x\n", sourceIPAddr)

		// Destination IP Address
		// destinationIPAddr := binary.BigEndian.Uint32(iPHeader[16:20])
		// fmt.Printf("Destination IP Address: %x\n", destinationIPAddr)

		// Protocol
		// iPProtocol := iPHeader[9]
		// fmt.Printf("Transport Protocol: %x\n", iPProtocol)

		// Transport Packet

		tcpHeader := make([]byte, 20)
		count, err = file.Read(tcpHeader)
		if err != nil {
			log.Fatal("Error reading TCP header")
		}

		sourcePort := binary.BigEndian.Uint16(tcpHeader[0:2])
		// destPort := binary.BigEndian.Uint16(tcpHeader[2:4])
		// fmt.Printf("Source Port: %d %d\n", sourcePort, destPort)

		dataOffset := uint16(tcpHeader[12] >> 4)

		seqNum := binary.BigEndian.Uint32(tcpHeader[4:8])
		fmt.Printf("Sequence Number: %d\n", seqNum)

		// dataSectionLen := totalIPLength - 20 - dataOffset*4
		// fmt.Printf("Packet Length: %d; IP Length: %d; Data Offset: %d; Data Section Length: %d\n", packetLength, totalIPLength, dataOffset, dataSectionLen)

		// read remaining part of IP Header
		ipOptions := make([]byte, dataOffset*4-20)
		count, err = file.Read(ipOptions)
		if err != nil {
			log.Fatal("Error reading IP Options")
		}

		// Parse rest of packet
		packetData := make([]byte, packetLength-6-6-2-20-dataOffset*4)
		count, err = file.Read(packetData)
		if err != nil {
			log.Fatal(err)
		}

		if sourcePort == 80 {
			// fmt.Printf("Packet %d Data: %x\n", packetNum, packetData)
			httpData[seqNum] = packetData
		}

		packetNum++
	}

	fmt.Println(httpData)
}
