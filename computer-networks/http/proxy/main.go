package main

import (
	"fmt"
	"syscall"
)

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println(err)
	}

	addr := &syscall.SockaddrInet4{Port: 8000}
	err = syscall.Bind(fd, addr)
	if err != nil {
		fmt.Println(err)
	}

	err = syscall.Listen(fd, 10)
	if err != nil {
		fmt.Println(err)
	}

	nfd, sa, err := syscall.Accept(fd)
	if err != nil {
		fmt.Println("accepting:", err)
	}

	fmt.Printf("%v %v\n", nfd, sa)

	buf := make([]byte, 100)
	n, from, err := syscall.Recvfrom(nfd, buf, 0)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Received:", n, from)

	err = syscall.Sendto(nfd, buf, 0, addr)
	if err != nil {
		fmt.Println(err)
	}

	err = syscall.Close(fd)
	if err != nil {
		fmt.Println(err)
	}

}
