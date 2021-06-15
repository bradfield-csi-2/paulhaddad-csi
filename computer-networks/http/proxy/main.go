package main

import (
	"bytes"
	"fmt"
	"syscall"
)

func closeSocket(fd int) {
	err := syscall.Close(fd)
	if err != nil {
		fmt.Println(err)
	}
}

var cachePaths = [][]byte{[]byte("website/")}
var cachedReqs = make(map[string][]byte)

func forwardReq(req []byte) []byte {
	sfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println(err)
	}

	servAddr := &syscall.SockaddrInet4{Port: 9000}
	err = syscall.Connect(sfd, servAddr)
	if err != nil {
		fmt.Println(err)
	}

	err = syscall.Sendto(sfd, req, 0, servAddr)
	if err != nil {
		fmt.Println(err)
	}

	rBuf := make([]byte, 1024)
	n, _, err := syscall.Recvfrom(sfd, rBuf, 0)
	if err != nil {
		fmt.Println(err)
	}

	err = syscall.Close(sfd)
	if err != nil {
		fmt.Println(err)
	}

	return rBuf[:n]
}

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println(err)
	}
	defer closeSocket(fd)

	proxyAddr := &syscall.SockaddrInet4{Port: 8000}
	err = syscall.Bind(fd, proxyAddr)
	if err != nil {
		fmt.Println(err)
	}

	err = syscall.Listen(fd, 10)
	if err != nil {
		fmt.Println(err)
	}

	var resp []byte
	for {
		nfd, _, err := syscall.Accept(fd)
		if err != nil {
			fmt.Println("accepting:", err)
		}

		buf := make([]byte, 1024)
		n, _, err := syscall.Recvfrom(nfd, buf, 0)
		if err != nil {
			fmt.Println(err)
		}

		reqMsg := bytes.Split(buf[:n], []byte("\n"))
		reqLine := bytes.Split(reqMsg[0], []byte(" "))

		path := reqLine[1][1:]
		fmt.Println("Path: ", path)

		searchCache := matchesCachePath(path)
		cachedRes, found := cachedReqs[string(path)]

		switch {
		case searchCache && found:
			fmt.Println("Found in cache")
			resp = cachedRes
		case searchCache && !found:
			fmt.Println("Need to cache")
			resp = forwardReq(buf[:n])
			cachedReqs[string(path)] = resp
		default:
			fmt.Println("Bypass cache")
			resp = forwardReq(buf[:n])
		}

		err = syscall.Sendto(nfd, resp, 0, proxyAddr)
		if err != nil {
			fmt.Println(err)
		}

		err = syscall.Close(nfd)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func matchesCachePath(path []byte) bool {
	for _, cp := range cachePaths {
		if bytes.HasPrefix(path, cp) {
			return true
		}
	}

	return false
}
