package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"syscall"
)

func main() {
	// open connection
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		fmt.Println("error opening socket connection")
	}

	// bind to port
	srcAddr := syscall.SockaddrInet4{}
	err = syscall.Bind(fd, &srcAddr)
	if err != nil {
		fmt.Println(err)
	}

	// build message
	var buf bytes.Buffer

	// header
	err = binary.Write(&buf, binary.BigEndian, uint16(16))
	if err != nil {
		fmt.Println(err)
	}

	err = binary.Write(&buf, binary.BigEndian, uint16(0b0000000100000000))
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(&buf, binary.BigEndian, uint16(1))
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(&buf, binary.BigEndian, uint16(0))
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(&buf, binary.BigEndian, uint16(0))
	if err != nil {
		fmt.Println(err)
	}
	err = binary.Write(&buf, binary.BigEndian, uint16(0))
	if err != nil {
		fmt.Println(err)
	}

	// question
	_, err = buf.Write([]byte{11, 'p', 'a', 'u', 'l', 'g', 'h', 'a', 'd', 'd', 'a', 'd', 3, 'c', 'o', 'm', 0})
	if err != nil {
		fmt.Println(err)
	}

	_, err = buf.Write([]byte{0, 1, 0, 1})
	if err != nil {
		fmt.Println(err)
	}

	// send to destination
	toAddr := syscall.SockaddrInet4{Port: 53, Addr: [4]byte{8, 8, 8, 8}}
	err = syscall.Sendto(fd, buf.Bytes(), 0, &toAddr)
	if err != nil {
		fmt.Println(err)
	}

	// receive message
	recvBuf := make([]byte, 250)
	_, _, err = syscall.Recvfrom(fd, recvBuf, 0)
	if err != nil {
		fmt.Println(err)
	}

	// decode response
	respHeader := recvBuf[0:12]

	// transID := binary.BigEndian.Uint16(respHeader[0:2])

	qr := recvBuf[3] >> 7
	if qr != 1 {
		fmt.Println("QR code should be a response (1)")
	}

	rcode := recvBuf[4] & 0b00001111
	if rcode != 0 {
		fmt.Printf("response code had an error condition: %d", rcode)
	}

	answerCount := binary.BigEndian.Uint16(respHeader[6:8])

	// questionSect := recvBuf[12 : 12+17+4]

	answerStart := 12 + 17 + 4
	answerLength := 16
	for i := 0; i < int(answerCount); i++ {
		var typeRec, classVal string

		answerSect := recvBuf[answerStart+answerLength*i : answerStart+answerLength*i+answerLength]

		// namePtr := binary.BigEndian.Uint16(answerSect[0:2]) & 0b0011111111111111

		// need to decode name encoding to string
		// name := recvBuf[namePtr : namePtr+17]

		// type
		typeRecVal := binary.BigEndian.Uint16(answerSect[2:4])

		switch typeRecVal {
		case 1:
			typeRec = "A"
		}

		// class
		rdataClassVal := binary.BigEndian.Uint16(answerSect[4:6])

		switch rdataClassVal {
		case 1:
			classVal = "IN"
		}

		// ttl
		ttl := binary.BigEndian.Uint32(answerSect[6:10])

		// data length
		rdLength := binary.BigEndian.Uint16(answerSect[10:12])

		// rdata
		rdata := answerSect[12 : 12+rdLength]

		fmt.Printf("%s\t%d\t%s\t%s\t%d\n", "paulghaddad.com", ttl, classVal, typeRec, rdata)
	}
}
