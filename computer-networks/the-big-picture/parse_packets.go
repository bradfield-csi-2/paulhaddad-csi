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
	globalHeaderLength            = 24
	ipV4EtherType          uint16 = 0x0008
	ipV6EtherType          uint16 = 0xDD86
	ipHeaderLength                = 20
	pcapGlobalHeaderLength        = 24
	packetHeaderLength            = 16
	ethernetHeaderLength          = 14
)

func parsePcapGlobalHeader(f *os.File) error {
	globalHeader := make([]byte, pcapGlobalHeaderLength)
	globalCount, err := f.Read(globalHeader)

	if err != nil {
		return fmt.Errorf("error reading global header")
	}

	if globalCount != globalHeaderLength {
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

func parseEthernetHeader(f *os.File) (uint16, error) {
	macDestination := make([]byte, 6)

	_, err := f.Read(macDestination)
	if err != nil {
		return 0, fmt.Errorf("error reading mac destination address")
	}

	macSource := make([]byte, 6)
	_, err = f.Read(macSource)
	if err != nil {
		return 0, fmt.Errorf("error reading mac source address")
	}

	etherTypeData := make([]byte, 2)
	_, err = f.Read(etherTypeData)
	if err != nil {
		return 0, fmt.Errorf("error reading EtherType field")
	}

	etherType := binary.LittleEndian.Uint16(etherTypeData)

	if etherType != ipV4EtherType && etherType != ipV6EtherType {
		return 0, fmt.Errorf("the IP datagram must be IPv4 or IPv6")
	}

	return etherType, nil
}

func parseIPHeader(f *os.File) error {
	ipHeader := make([]byte, ipHeaderLength)

	_, err := f.Read(ipHeader)
	if err != nil {
		return fmt.Errorf("error reading IP header")
	}

	return nil
}

func parseTCPHeader(f *os.File) (uint16, uint32, uint16, error) {
	tcpHeader := make([]byte, 20)
	_, err := f.Read(tcpHeader)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("Error reading TCP header")
	}

	sourcePort := binary.BigEndian.Uint16(tcpHeader[0:2])
	dataOffset := uint16(tcpHeader[12] >> 4)
	seqNum := binary.BigEndian.Uint32(tcpHeader[4:8])

	ipOptions := make([]byte, dataOffset*4-ipHeaderLength)
	_, err = f.Read(ipOptions)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("Error reading IP Options")
	}

	return sourcePort, seqNum, dataOffset, nil
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
	httpComponents := bytes.Split(binaryString, []byte{13, 10})
	httpBody := httpComponents[len(httpComponents)-1]

	f, err := os.Create("image.jpg")
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
	file, err := os.Open(filename)
	if err != nil {
		return err
	}

	parsePcapGlobalHeader(file)

	var packetNum int
	var firstEtherType uint16
	httpData := make(map[int][]byte)

	for {
		packetLength, err := parsePacketHeader(file)
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		etherType, err := parseEthernetHeader(file)
		if err != nil {
			break
		}

		if packetNum == 0 {
			firstEtherType = etherType
		}

		if packetNum > 0 && etherType != firstEtherType {
			break
		}

		err = parseIPHeader(file)
		if err != nil {
			break
		}

		sourcePort, seqNum, dataOffset, err := parseTCPHeader(file)
		if err != nil {
			break
		}

		packetData := make([]byte, packetLength-ethernetHeaderLength-ipHeaderLength-dataOffset*4)
		_, err = file.Read(packetData)
		if err != nil {
			break
		}

		if sourcePort == 80 {
			httpData[int(seqNum)] = packetData
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
	err := parsePackets("net.cap")

	if err != nil {
		fmt.Printf("error: %s", err)
	}
}
