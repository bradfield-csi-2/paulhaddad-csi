package main

import (
	"bytes"
	"fmt"
	"log"
	"syscall"
)

const (
	proxyPort  = 8000
	serverPort = 9000
)

var cachePaths = [][]byte{[]byte("website/")}
var cachedReqs = make(map[string][]byte)

func closeSocket(fd int) {
	err := syscall.Close(fd)
	if err != nil {
		fmt.Println(err)
	}
}

func parsePath(buf []byte) []byte {
	reqMsg := bytes.Split(buf, []byte("\n"))
	reqLine := bytes.Split(reqMsg[0], []byte(" "))
	return reqLine[1][1:]
}

func forwardReq(req []byte) ([]byte, error) {
	var rBuf []byte = make([]byte, 1024)

	sfd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return nil, fmt.Errorf("error creating socket: %s", err)
	}

	defer closeSocket(sfd)

	servAddr := &syscall.SockaddrInet4{Port: serverPort}
	err = syscall.Connect(sfd, servAddr)
	if err != nil {
		return nil, fmt.Errorf("error connecting to socket: %s", err)
	}

	err = syscall.Sendto(sfd, req, 0, servAddr)
	if err != nil {
		return nil, fmt.Errorf("error sending to proxy: %s", err)
	}

	n, _, err := syscall.Recvfrom(sfd, rBuf, 0)
	if err != nil {
		return nil, fmt.Errorf("error receiving from proxy: %s", err)
	}

	return rBuf[:n], nil
}

func matchesCachePath(path []byte) bool {
	for _, cp := range cachePaths {
		if bytes.HasPrefix(path, cp) {
			return true
		}
	}

	return false
}

func createProxyResp(resp []byte) []byte {
	proxyRespHdr := bytes.Split(resp, []byte("\r\n\r\n"))[0]
	proxyRespData := bytes.Split(resp, []byte("\r\n\r\n"))[1]

	updatedHeader := make([]byte, 0)
	updatedHeader = append(updatedHeader, proxyRespHdr...)

	// add keep-alive header for persistant connection
	updatedHeader = append(updatedHeader, []byte("\r\nConnection: Keep-Alive\r\n\r\n")...)

	proxiedResp := make([]byte, 0)
	proxiedResp = append(proxiedResp, updatedHeader...)
	proxiedResp = append(proxiedResp, proxyRespData...)

	return proxiedResp
}

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		fmt.Println("error connecting to socket:", err)
	}
	defer closeSocket(fd)

	proxyAddr := &syscall.SockaddrInet4{Port: proxyPort}
	err = syscall.Bind(fd, proxyAddr)
	if err != nil {
		fmt.Println("error binding to socket:", err)
	}

	err = syscall.Listen(fd, 10)
	if err != nil {
		fmt.Println(err)
	}

	for {
		var proxyResp []byte

		nfd, _, err := syscall.Accept(fd)
		if err != nil {
			log.Fatal("error accepting:", err)
		}

		defer closeSocket(nfd)

		buf := make([]byte, 1024)
		n, _, err := syscall.Recvfrom(nfd, buf, 0)
		if err != nil {
			log.Fatal("error receiving from client:", err)
		}

		path := parsePath(buf[:n])

		searchCache := matchesCachePath(path)
		cachedRes, found := cachedReqs[string(path)]

		switch {
		case searchCache && found:
			log.Println("Found in cache")
			proxyResp = cachedRes
		case searchCache && !found:
			log.Println("Need to cache")
			serverResp, err := forwardReq(buf[:n])
			if err != nil {
				log.Fatal("error encountered forwarding request")
			}
			proxyResp = createProxyResp(serverResp)
			cachedReqs[string(path)] = proxyResp
		default:
			log.Println("Bypass cache")
			serverResp, err := forwardReq(buf[:n])
			if err != nil {
				log.Fatal("error encountered forwarding request")
			}
			proxyResp = createProxyResp(serverResp)
		}

		err = syscall.Sendto(nfd, proxyResp, 0, proxyAddr)
		if err != nil {
			log.Fatal("error sending to client:", err)
		}
	}
}
