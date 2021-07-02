package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"reflect"
	"syscall"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("traceroute requires a host")
		os.Exit(1)
	}
	hostname := os.Args[1]

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

	addrs, err := net.LookupIP(hostname)
	if err != nil {
		fmt.Printf("error looking up host: %s\n", err)
		os.Exit(1)
	}

	addr := addrs[0]
	fmt.Printf("traceroute to %s (%v), 64 hops max, 52 byte packets\n", hostname, addr)
	host := [4]byte{addr[12], addr[13], addr[14], addr[15]}

	var seqNum uint16 = 1
	var icmpType byte
	port := 33434
	for {
		times := make([]float64, 3)
		sourceIPs := make([][]byte, 3)

		for i := 0; i < 3; i++ {
			start := time.Now()
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

			// determine RTT
			end := time.Now()
			times[i] = float64(end.Sub(start).Microseconds()) / 1000.0

			// determine sourceIP
			headerLen := uint16(resp[0] & 0b00001111 * 4)
			ipLength := binary.LittleEndian.Uint16(resp[2:4])
			totalIPLen := headerLen + ipLength
			ipPacket := resp[0:totalIPLen]
			sourceIPs[i] = ipPacket[12:16]

			// determine ICMP type
			icmpFrame := ipPacket[headerLen:]
			icmpType = icmpFrame[0]
		}

		// determine if all the IPs are the same
		sameIPs := true
		prevIP := sourceIPs[0]
		for i := 1; i < 3; i++ {
			if !reflect.DeepEqual(prevIP, sourceIPs[i]) {
				sameIPs = false
				break
			}
			prevIP = sourceIPs[i]
		}

		if sameIPs {
			fmt.Printf(" %d  %d.%d.%d.%d  %.3f ms  %.3f ms  %.3f ms\n",
				seqNum, sourceIPs[0][0], sourceIPs[0][1], sourceIPs[0][2],
				sourceIPs[0][3], times[0], times[1], times[2])
		} else {
			fmt.Printf(" %d  %d.%d.%d.%d %.3f ms\n", seqNum, sourceIPs[0][0],
				sourceIPs[0][1], sourceIPs[0][2], sourceIPs[0][3], times[0])

			for i := 1; i < 3; i++ {
				fmt.Printf("    %d.%d.%d.%d %.3f ms\n", sourceIPs[i][0],
					sourceIPs[i][1], sourceIPs[i][2], sourceIPs[i][3], times[i])
			}
		}

		// destination unreachable?
		if icmpType == 3 {
			break
		}

		seqNum++
		port++
	}
}
