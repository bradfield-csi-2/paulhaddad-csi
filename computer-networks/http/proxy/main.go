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

	proxyAddr := &syscall.SockaddrInet4{Port: 8000}
	err = syscall.Bind(fd, proxyAddr)
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

	buf := make([]byte, 1024)
	n, from, err := syscall.Recvfrom(nfd, buf, 0)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Received:", n, from, buf)

	// connect to backing server and send message
	sfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println(err)
	}

	servAddr := &syscall.SockaddrInet4{Port: 9000}
	err = syscall.Connect(sfd, servAddr)
	if err != nil {
		fmt.Println(err)
	}

	err = syscall.Sendto(sfd, buf, 0, servAddr)
	if err != nil {
		fmt.Println(err)
	}

	rBuf := make([]byte, 1024)
	n, from, err = syscall.Recvfrom(sfd, rBuf, 0)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Received from server: %s\n\n", rBuf)

	err = syscall.Sendto(nfd, rBuf, 0, proxyAddr)
	if err != nil {
		fmt.Println(err)
	}

	err = syscall.Close(sfd)
	if err != nil {
		fmt.Println(err)
	}

	err = syscall.Close(nfd)
	if err != nil {
		fmt.Println(err)
	}

	err = syscall.Close(fd)
	if err != nil {
		fmt.Println(err)
	}

}
