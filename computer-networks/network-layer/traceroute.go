package main

import (
	"encoding/binary"
	"fmt"
	"os"
	"syscall"
)

func main() {
	sender, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		fmt.Printf("error opening sender socket: %s\n", err)
		os.Exit(1)
	}

	receiver, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	if err != nil {
		fmt.Printf("error opening receiver socket: %s\n", err)
		os.Exit(1)
	}

	// addrs, err := net.LookupIP("google.com")
	// if err != nil {
	// 	fmt.Printf("error looking up host: %s\n", err)
	// 	os.Exit(1)
	// }

	// hardcoding bradfield's IP in for now
	host := [4]byte{104, 21, 76, 199}

	var idenNum uint16 = 1
	var seqNum uint16 = 9
	port := 33434

	for {
		fmt.Printf("Sending hop: %d\n", seqNum)
		err = syscall.SetsockoptInt(sender, syscall.IPPROTO_IP, syscall.IP_TTL, int(seqNum))
		if err != nil {
			fmt.Printf("error socket option: %s\n", err)
			os.Exit(1)
		}

		sockAddr := syscall.SockaddrInet4{Addr: host, Port: port}
		reqMsg := make([]byte, 24)
		err = syscall.Sendto(sender, reqMsg, 0, &sockAddr)
		if err != nil {
			fmt.Printf("error sending: %s\n", err)
			os.Exit(1)
		}

		resp := make([]byte, 1024)
		_, _, err = syscall.Recvfrom(receiver, resp, 0)
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

		fmt.Printf("Type: %d; Code: %d\n", icmpType, icmpCode)

		origReq := icmpFrame[8:]
		intIPheaderLen := uint16(origReq[0] & 0b00001111 * 4)
		intIPLength := binary.LittleEndian.Uint16(origReq[2:4])

		fmt.Printf("Header Length: %d ICMP Data Length: %d\n", intIPheaderLen, intIPLength)

		icmpErrMsg := origReq[intIPheaderLen : intIPheaderLen+intIPLength]
		fmt.Printf("ICMP Error Message: % x\n", icmpErrMsg)

		respIdenNum := binary.BigEndian.Uint16(icmpErrMsg[4:6])
		respSeqNum := binary.BigEndian.Uint16(icmpErrMsg[6:8])
		fmt.Printf("Identification: %d Response Identifier: %d\n", idenNum, respIdenNum)
		fmt.Printf("Seq Num: %d Response Sequence Num: %d\n", seqNum, respSeqNum)

		if icmpType == 3 {
			break
		}

		seqNum++
		port++
	}
}
