package net

import (
	"errors"
	"log"
	"net"
	"syscall"
)

// Htons converts host data to network byte order.
func Htons(v uint16) uint16 {
	return htons(v)
}

// little-endian to network byte order
func htons(v uint16) uint16 {
	return (v << 8) + (v >> 8)
}

// NewRawSocket returns file discriptor
func NewRawSocket() int {
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_IP)))
	if err != nil {
		log.Fatalf("create udp sendfd err : %v\n", err)
	}
	return fd
}

func GetDeviceIpAddr(device string) (string, error) {
	inf, err := net.InterfaceByName(device)
	if err != nil {
		return "", err
	}
	addrs, err := inf.Addrs()
	if err != nil {
		return "", err
	}
	if len(addrs) == 0 {
		return "", errors.New("addr not found")
	}
	return addrs[0].String(), nil
}
