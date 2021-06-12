package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"syscall"
)

type DNSMsg struct {
	txnID    uint16
	msgType  uint8
	rCode    uint8
	qstCount uint16
	ansCount uint16
	nsCount  uint16
	arCount  uint16
	question Question
	answers  []ResourceRecord
}

type Question struct {
	qName  []byte
	qType  uint16
	qClass uint16
}

type ResourceRecord struct {
	name     []byte
	rrType   string
	rrClass  string
	ttl      uint32
	rdLength uint16
	rData    []byte
}

func openSocket() (int, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)

	return fd, err
}

func sendQuery(fd int, addr [4]byte, port int, query []byte) error {
	toAddr := syscall.SockaddrInet4{Port: port, Addr: addr}

	return syscall.Sendto(fd, query, 0, &toAddr)
}

func recvResponse(fd int, buf []byte) error {
	_, _, err := syscall.Recvfrom(fd, buf, 0)

	return err
}

func encodeMsg() []byte {
	b := make([]byte, 12)

	// header
	binary.BigEndian.PutUint16(b[0:], 0x0010)
	binary.BigEndian.PutUint16(b[2:], uint16(0b0000000100000000))
	binary.BigEndian.PutUint16(b[4:], 0x0001)
	binary.BigEndian.PutUint16(b[6:], 0x0000)
	binary.BigEndian.PutUint16(b[8:], 0x0000)
	binary.BigEndian.PutUint16(b[10:], 0x0000)

	// question
	b = append(b, []byte{11, 'p', 'a', 'u', 'l', 'g', 'h', 'a', 'd', 'd', 'a', 'd', 3, 'c', 'o', 'm', 0}...)
	b = append(b, []byte{0, 1, 0, 1}...)

	return b
}

func decodeResp(b []byte) *DNSMsg {
	var msg DNSMsg

	respHeader := b[0:12]

	msg.txnID = binary.BigEndian.Uint16(respHeader[0:2])
	msg.msgType = b[3] >> 7
	// if qr != 1 {
	// 	fmt.Println("QR code should be a response (1)")
	// }

	msg.rCode = b[4] & 0b00001111
	// if rcode != 0 {
	// 	fmt.Printf("response code had an error condition: %d", rcode)
	// }

	msg.ansCount = binary.BigEndian.Uint16(respHeader[6:8])

	// questionSect := b[12 : 12+17+4]

	answerStart := 12 + 17 + 4
	answerLength := 16
	for i := 0; i < int(msg.ansCount); i++ {
		var rr ResourceRecord

		answerSect := b[answerStart+answerLength*i : answerStart+answerLength*i+answerLength]

		// namePtr := binary.BigEndian.Uint16(answerSect[0:2]) & 0b0011111111111111

		// need to decode name encoding to string
		// name := recvBuf[namePtr : namePtr+17]

		typeRecVal := binary.BigEndian.Uint16(answerSect[2:4])

		switch typeRecVal {
		case 1:
			rr.rrType = "A"
		}

		rdataClassVal := binary.BigEndian.Uint16(answerSect[4:6])

		switch rdataClassVal {
		case 1:
			rr.rrClass = "IN"
		}

		rr.ttl = binary.BigEndian.Uint32(answerSect[6:10])
		rr.rdLength = binary.BigEndian.Uint16(answerSect[10:12])
		rr.rData = answerSect[12 : 12+rr.rdLength]
		msg.answers = append(msg.answers, rr)
	}

	return &msg
}

func printResp(records *[]ResourceRecord) {
	for _, rr := range *records {
		fmt.Printf("%s\t%d\t%s\t%s\t%d\n", "paulghaddad.com", rr.ttl, rr.rrClass, rr.rrType, rr.rData)
	}
}

func main() {
	// TODO: Take in command line arg
	fd, err := openSocket()
	if err != nil {
		fmt.Printf("error opening socket connection: %s\n", err)
		os.Exit(1)
	}

	queryMsg := encodeMsg()

	err = sendQuery(fd, [4]byte{8, 8, 8, 8}, 53, queryMsg)
	if err != nil {
		fmt.Printf("error sending query: %s\n", err)
		os.Exit(1)
	}

	recvBuf := make([]byte, 250)
	err = recvResponse(fd, recvBuf)
	if err != nil {
		fmt.Printf("error receiving response: %s\n", err)
		os.Exit(1)
	}

	resp := decodeResp(recvBuf)

	printResp(&resp.answers)

	// need to close the socket
}
