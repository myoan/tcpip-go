package main

import (
	"flag"
	"fmt"
	"net"
	"syscall"

	"github.com/BurntSushi/toml"
	mynet "github.com/myoan/tcpip-go"
)

type Config struct {
	IpTTL  int    `toml:"ip-ttl"`
	Device string `toml:"device"`
	VMac   string `toml:"vmac"`
	VIP    string `toml:"vip"`
	VMask  string `toml:"vmask"`
}

func main() {
	var confPath string
	flag.StringVar(&confPath, "c", "", "config file path")
	flag.Parse()

	var conf Config
	_, err := toml.DecodeFile(confPath, &conf)
	if err != nil {
		panic(err)
	}

	fmt.Printf("ip-ttl: %d\n", conf.IpTTL)
	fmt.Printf("device: %s\n", conf.Device)
	fmt.Printf("VMac: %s\n", conf.VMac)
	fmt.Printf("VIP: %s\n", conf.VIP)
	fmt.Printf("VMask: %s\n", conf.VMask)

	inf, err := net.InterfaceByName("enp0s5")
	if err != nil {
		panic(err)
	}

	// addr, err := mynet.GetDeviceIpAddr("enp0s5")
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(addr)

	fmt.Printf("mac: %v\n", inf.HardwareAddr)

	sockfd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(mynet.Htons(syscall.ETH_P_ALL)))
	if err != nil {
		panic(err)
	}
	defer syscall.Close(sockfd)

	layer := syscall.SockaddrLinklayer{
		Protocol: syscall.ETH_P_ARP,
		Ifindex:  inf.Index,
		Hatype:   syscall.ARPHRD_ETHER,
	}

	// destination MAC addr
	pkt := []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff}
	// source MAC addr
	pkt = append(pkt, []byte(inf.HardwareAddr)...)
	// type
	pkt = append(pkt, []byte{0x08, 0x06}...)

	arp := mynet.NewArpRequest(
		[6]byte(inf.HardwareAddr),
		[6]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		[4]byte{0x0a, 0xd3, 0x37, 0x17},
		[4]byte{0x0a, 0xd3, 0x37, 0x01},
	)
	pkt = append(pkt, arp.ToPacket()...)
	fmt.Printf("pkt: %v\n", pkt)

	err = syscall.Sendto(sockfd, pkt, 0, &layer)
	if err != nil {
		panic(err)
	}

	for {
		rcvBuf := make([]byte, 80)
		_, _, err = syscall.Recvfrom(sockfd, rcvBuf, 0)
		if err != nil {
			panic(err)
		}
		if rcvBuf[12] == 0x08 && rcvBuf[13] == 0x06 {
			var frame mynet.EthernetFrame
			fmt.Printf("receive: %v\n", rcvBuf)
			err := mynet.UnmarshallEtherFrame(&frame, rcvBuf)
			if err != nil {
				panic(err)
			}
			fmt.Printf("dst: %s\n", frame.DstMAC())
			fmt.Printf("src: %s\n", frame.SrcMAC())

			var arp mynet.ArpPacket
			err = mynet.UnmarshallArpPacket(&arp, frame.Data)
			if err != nil {
				panic(err)
			}
			fmt.Printf("hwtype: %02X, prtype: %02X, twlen: %02X, prtlen: %02X\n", arp.HardwareType, arp.ProtocolType, arp.HardwareLen, arp.ProtocolLen)
			fmt.Printf("op: %02X\n", arp.Operation)
			fmt.Printf("sender: %s\n", arp.SenderHWAddress())
			fmt.Printf("target: %s\n", arp.TargetHWAddress())
			break
		}
	}
}
