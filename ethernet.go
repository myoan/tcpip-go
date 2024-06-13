package net

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
)

type Ethernet struct {
	fd     int
	packet EthernetPacket
}

type EthernetPacket struct {
	Preamble [8]byte
	Frame    EthernetFrame
	FCS      [4]byte
}

const (
	ETHERNET_TYPE_IP4 uint16 = iota
	ETHERNET_TYPE_ARP
)

var ethernetTypeMap = map[uint16]uint16{
	ETHERNET_TYPE_IP4: 0x0000,
	ETHERNET_TYPE_ARP: 0x0806,
}

type EthernetFrame struct {
	dstMAC [6]byte
	srcMAC [6]byte
	Type   uint16
	Data   []byte
}

func UnmarshallEtherFrame(frame *EthernetFrame, data []byte) error {
	if len(data) < 14 {
		return errors.New("data too short")
	}
	copy(frame.dstMAC[:], data[0:6])
	copy(frame.srcMAC[:], data[6:12])
	typeUint := uint16(binary.BigEndian.Uint16(data[12:14]))
	for k, v := range ethernetTypeMap {
		if v == typeUint {
			frame.Type = k
			frame.Data = data[14:]
		}
	}
	if frame.Type == 0 {
		return fmt.Errorf("unsupported type: %d", typeUint)
	}
	return nil
}

func byteToMAC(v [6]byte) string {
	return fmt.Sprintf("%02X:%02X:%02X:%02X:%02X:%02X", v[0], v[1], v[2], v[3], v[4], v[5])
}

func (f *EthernetFrame) DstMAC() string {
	return byteToMAC(f.dstMAC)
}

func (f *EthernetFrame) SrcMAC() string {
	return byteToMAC(f.srcMAC)
}

type ArpPacket struct {
	HardwareType          [2]byte
	ProtocolType          [2]byte
	HardwareLen           [1]byte
	ProtocolLen           [1]byte
	Operation             [2]byte
	senderHWAddress       [6]byte
	senderProtocolAddress [4]byte
	targetHWAddress       [6]byte
	targetProtocolAddress [4]byte
}

func UnmarshallArpPacket(pkt *ArpPacket, data []byte) error {
	copy(pkt.HardwareType[:], data[0:2])
	copy(pkt.ProtocolType[:], data[2:4])
	copy(pkt.HardwareLen[:], data[4:5])
	copy(pkt.ProtocolLen[:], data[5:6])
	copy(pkt.Operation[:], data[6:8])
	copy(pkt.senderHWAddress[:], data[8:14])
	copy(pkt.senderProtocolAddress[:], data[14:18])
	copy(pkt.targetProtocolAddress[:], data[18:24])
	copy(pkt.targetProtocolAddress[:], data[24:28])
	return nil
}

func NewArpRequest(srcAddr, dstAddr [6]byte, srcIP, dstIP [4]byte) *ArpPacket {
	return &ArpPacket{
		HardwareType:          [2]byte{0x00, 0x01},
		ProtocolType:          [2]byte{0x08, 0x00},
		HardwareLen:           [1]byte{0x06},
		ProtocolLen:           [1]byte{0x04},
		Operation:             [2]byte{0x00, 0x01},
		senderHWAddress:       srcAddr,
		senderProtocolAddress: srcIP,
		targetHWAddress:       dstAddr,
		targetProtocolAddress: dstIP,
	}
}

func (a *ArpPacket) SenderHWAddress() string {
	return byteToMAC(a.senderHWAddress)
}

func (a *ArpPacket) TargetHWAddress() string {
	return byteToMAC(a.targetHWAddress)
}

func (a *ArpPacket) ToPacket() []byte {
	ret := make([]byte, 0)
	buf := bytes.NewBuffer(ret)
	buf.Write(a.HardwareType[:])
	buf.Write(a.ProtocolType[:])
	buf.Write(a.HardwareLen[:])
	buf.Write(a.ProtocolLen[:])
	buf.Write(a.Operation[:])
	buf.Write(a.senderHWAddress[:])
	buf.Write(a.senderProtocolAddress[:])
	buf.Write(a.targetHWAddress[:])
	buf.Write(a.targetProtocolAddress[:])
	return buf.Bytes()
}

var (
	EthernetTypeAddressResolutionProtocol = 0x0806
)

func NewEthernet(fd int) *Ethernet {
	return &Ethernet{
		fd: fd,
	}
}

func NewEthernetPacket(dst, src string, typ int) EthernetPacket {
	return EthernetPacket{
		Preamble: [8]byte{0, 0, 0, 0, 0, 0, 1, 0},
		Frame: EthernetFrame{
			dstMAC: [6]byte{},
			srcMAC: [6]byte{},
			Type:   0,
			Data:   []byte{},
		},
		FCS: [4]byte{0, 0, 0, 0},
	}
}

func (e *Ethernet) Send() error {
	return nil
}
