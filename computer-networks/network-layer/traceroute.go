package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"syscall"
)

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		fmt.Printf("error opening socket: %s\n", err)
		os.Exit(1)
	}

	err = syscall.SetsockoptInt(fd, syscall.IPPROTO_IP, syscall.IP_TTL, 1)
	if err != nil {
		fmt.Printf("error socket option: %s\n", err)
		os.Exit(1)
	}

	// addrs, err := net.LookupIP("google.com")
	// if err != nil {
	// 	fmt.Printf("error looking up host: %s\n", err)
	// 	os.Exit(1)
	// }

	// hardcoding google's IP in for now
	host := [4]byte{172, 67, 200, 47}

	// encode message
	msg := make([]byte, 10)
	binary.BigEndian.PutUint16(msg[0:], 0x0800)
	binary.BigEndian.PutUint16(msg[2:], 0xe7ff)
	binary.BigEndian.PutUint16(msg[4:], 0x1000) // Identifier Num
	binary.BigEndian.PutUint16(msg[6:], 0x0000) // Sequence Num
	binary.BigEndian.PutUint16(msg[8:], 0x0000) // add arbitrary data

	sockAddr := syscall.SockaddrInet4{Addr: host}
	err = syscall.Sendto(fd, msg, 0, &sockAddr)
	if err != nil {
		fmt.Printf("error sending: %s\n", err)
		os.Exit(1)
	}

	resp := make([]byte, 1024)
	_, _, err = syscall.Recvfrom(fd, resp, 0)
	if err != nil {
		fmt.Printf("error receiving: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Header: % x\n", resp[0:20])
	headerLen := uint16(resp[0] & 0b00001111 * 4)
	ipLength := binary.LittleEndian.Uint16(resp[2:4])
	totalIPLen := headerLen + ipLength

	ipPacket := resp[0:totalIPLen]
	fmt.Printf("IP Packet: %x\n", ipPacket)

	sourceIP := ipPacket[12:16]
	fmt.Printf("Source IP: %d.%d.%d.%d\n", sourceIP[0], sourceIP[1], sourceIP[2], sourceIP[3])

	icmpFrame := ipPacket[headerLen:]
	icmpType := icmpFrame[0]
	icmpCode := icmpFrame[1]

	// assert type == 11 and code == 0
	fmt.Printf("Type: %d; Code: %d\n", icmpType, icmpCode)

	origReq := icmpFrame[8:]
	intIPheaderLen := uint16(origReq[0] & 0b00001111 * 4)
	intIPLength := binary.LittleEndian.Uint16(origReq[2:4])

	fmt.Printf("Header Length: %d ICMP Data Length: %d\n", intIPheaderLen, intIPLength)

	icmpErrMsg := origReq[intIPheaderLen : intIPheaderLen+intIPLength]
	fmt.Printf("ICMP Error Message: % x\n", icmpErrMsg)

	respIdenNum := icmpErrMsg[4:6]
	respSeqNum := icmpErrMsg[6:8]

	fmt.Printf("Response Identifier: %x\n", respIdenNum)
	fmt.Printf("Response Sequence Num: %x\n", respSeqNum)
}
