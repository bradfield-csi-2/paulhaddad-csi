package main

import (
	"log"
	"syscall"
)

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_DGRAM, 0)
	if err != nil {
		log.Fatal("error creating socket")
	}

	to := syscall.SockaddrInet4{Port: 58604, Addr: [4]byte{0, 0, 0, 0}}
	err = syscall.Sendto(fd, []byte("hello"), 0, &to)
	if err != nil {
		log.Fatal("error sending")
	}

	// wait for ACK or NACK
	for {
		recvData := make([]byte, 1)
		_, _, err = syscall.Recvfrom(fd, recvData, 0)
		if err != nil {
			log.Fatal("error receiving")
		}
		ack := recvData[0]

		if ack == 0 {
			break
		}

		// resend packet
		log.Println("resending corrupted packet")
		err = syscall.Sendto(fd, []byte("hello"), 0, &to)
		if err != nil {
			log.Fatal("error sending")
		}
	}

	log.Println("Received packet successfully")

	// if corrupt: resend

}

// TODO:
// 1. Send packet to proxy - DONE
// 2. Receive packet from proxy - DONE
// 3. Handle dropped packets
// 4. Handle corrupted packets
