package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"sort"
)

const (
	// PCAP
	globalHeaderLength     = 24
	pcapGlobalHeaderLength = 24
	packetHeaderLength     = 16

	// Ethernet
	ethernetHeaderLength        = 14
	ipV4EtherType        uint16 = 0x0008
	ipV6EtherType        uint16 = 0xDD86

	// IP
	minIPHeaderLength = 20
	httpPort          = 80

	// TCP
	minTCPHeaderLength = 20

	inputFile  = "net.cap"
	outputFile = "image.jpg"
)

type ethernetHeaderData struct {
	macDestination uint64
	macSource      uint64
	etherType      uint16
}

type ipHeaderData struct {
	version           uint8
	ipHeaderLength    uint8
	totalLength       uint16
	sourceIPAddr      uint32
	destinationIPAddr uint32
}

type tcpHeaderData struct {
	sourcePort uint16
	destPort   uint16
	seqNum     uint32
	dataOffset uint16
}

func parsePcapGlobalHeader(f *os.File) error {
	globalHeader := make([]byte, pcapGlobalHeaderLength)
	count, err := f.Read(globalHeader)

	if err != nil {
		return fmt.Errorf("error reading global header")
	}

	if count != globalHeaderLength {
		return fmt.Errorf("pcap global header not correct length")
	}

	return nil
}

func parsePacketHeader(f *os.File) (uint16, error) {
	packetHeader := make([]byte, packetHeaderLength)

	count, err := f.Read(packetHeader)
	if err == io.EOF {
		return 0, err
	}

	if err != nil {
		return 0, fmt.Errorf("error parsing packet header")
	}

	if count != packetHeaderLength {
		return 0, fmt.Errorf("packet header not correct length")
	}

	packetLength := binary.LittleEndian.Uint16(packetHeader[8:11])
	untruncatedPacketLength := binary.LittleEndian.Uint16(packetHeader[12:15])

	if packetLength != untruncatedPacketLength {
		return 0, fmt.Errorf("packet is truncated")
	}

	return packetLength, nil
}

func parseEthernetHeader(f *os.File) (ethernetHeaderData, error) {
	var data ethernetHeaderData

	macDestinationData := make([]byte, 6)
	_, err := f.Read(macDestinationData)
	if err != nil {
		return data, fmt.Errorf("error reading mac destination address")
	}

	macDestination, _ := binary.Uvarint(macDestinationData)
	data.macDestination = macDestination

	macSourceData := make([]byte, 6)
	_, err = f.Read(macSourceData)
	if err != nil {
		return data, fmt.Errorf("error reading mac source address")
	}

	macSource, _ := binary.Uvarint(macSourceData)
	data.macSource = macSource

	etherTypeData := make([]byte, 2)
	_, err = f.Read(etherTypeData)
	if err != nil {
		return data, fmt.Errorf("error reading EtherType field")
	}

	etherType := binary.LittleEndian.Uint16(etherTypeData)

	if etherType != ipV4EtherType && etherType != ipV6EtherType {
		return data, fmt.Errorf("the IP datagram must be IPv4 or IPv6")
	}

	return data, nil
}

func parseIPHeader(f *os.File) (ipHeaderData, error) {
	var data ipHeaderData

	ipHeader := make([]byte, minIPHeaderLength)

	_, err := f.Read(ipHeader)
	if err != nil {
		return data, fmt.Errorf("error reading IP header")
	}

	data.version = uint8(ipHeader[0] >> 4)
	data.ipHeaderLength = uint8(ipHeader[0]&15) * 4
	data.totalLength = binary.BigEndian.Uint16(ipHeader[2:4])
	data.sourceIPAddr = binary.BigEndian.Uint32(ipHeader[12:16])
	data.destinationIPAddr = binary.BigEndian.Uint32(ipHeader[16:20])

	if data.ipHeaderLength > minIPHeaderLength {
		ipHeaderOptions := make([]byte, data.ipHeaderLength-minIPHeaderLength)
		_, err := f.Read(ipHeaderOptions)
		if err != nil {
			return data, fmt.Errorf("error reading IP header options")
		}
	}

	return data, nil
}

func parseTCPHeader(f *os.File) (tcpHeaderData, error) {
	var data tcpHeaderData

	tcpHeader := make([]byte, minTCPHeaderLength)
	_, err := f.Read(tcpHeader)
	if err != nil {
		return data, fmt.Errorf("Error reading TCP header")
	}

	data.sourcePort = binary.BigEndian.Uint16(tcpHeader[0:2])
	data.destPort = binary.BigEndian.Uint16(tcpHeader[2:4])
	data.dataOffset = uint16(tcpHeader[12] >> 4)
	data.seqNum = binary.BigEndian.Uint32(tcpHeader[4:8])

	tcpOptions := make([]byte, data.dataOffset*4-minTCPHeaderLength)
	_, err = f.Read(tcpOptions)
	if err != nil {
		return data, fmt.Errorf("Error reading IP Options")
	}

	return data, nil
}

func sortSeqNums(httpData map[int][]byte) []int {
	var seqNums []int

	for k := range httpData {
		seqNums = append(seqNums, k)
	}
	sort.Ints(seqNums)

	return seqNums
}

func writeHTTPData(seqNums []int, httpData map[int][]byte) error {
	var b bytes.Buffer

	for _, num := range seqNums {
		data := httpData[num]
		b.Write(data)
	}

	binaryString := b.Bytes()

	// split on CR-LF characters and get the http body
	httpComponents := bytes.Split(binaryString, []byte{13, 10})
	httpBody := httpComponents[len(httpComponents)-1]

	f, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("error creating file")
	}

	_, err = f.Write(httpBody)
	if err != nil {
		return fmt.Errorf("error writing to file")
	}

	return nil
}

func parsePackets(filename string) error {
	var packetNum int
	var firstEtherType uint16
	httpData := make(map[int][]byte)

	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	err = parsePcapGlobalHeader(file)
	if err != nil {
		return err
	}

	for {
		packetLength, err := parsePacketHeader(file)
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		ethernetHeader, err := parseEthernetHeader(file)
		if err != nil {
			break
		}

		// check that EtherType fields are the same
		if packetNum == 0 {
			firstEtherType = ethernetHeader.etherType
		}
		if packetNum > 0 && ethernetHeader.etherType != firstEtherType {
			break
		}

		ipHeader, err := parseIPHeader(file)
		if err != nil {
			break
		}

		tcpHeader, err := parseTCPHeader(file)
		if err != nil {
			break
		}

		// read tcp data
		tcpDataLength := packetLength - ethernetHeaderLength - uint16(ipHeader.ipHeaderLength) - tcpHeader.dataOffset*4
		tcpData := make([]byte, tcpDataLength)
		_, err = file.Read(tcpData)
		if err != nil {
			break
		}

		// filter HTTP responses from image server
		if tcpHeader.sourcePort == httpPort {
			httpData[int(tcpHeader.seqNum)] = tcpData
		}

		packetNum++
	}

	seqNums := sortSeqNums(httpData)

	err = writeHTTPData(seqNums, httpData)
	if err != nil {
		return fmt.Errorf("error writing http data to file")
	}

	return nil
}

func main() {
	err := parsePackets(inputFile)

	if err != nil {
		fmt.Printf("error: %s", err)
	}
}
