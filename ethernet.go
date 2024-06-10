package net

import "bytes"

type Ethernet struct {
	fd     int
	packet EthernetPacket
}

type EthernetPacket struct {
	Preamble [8]byte
	Frame    EthernetFrame
	FCS      [4]byte
}

type EthernetFrame struct {
	DstMAC [6]byte
	SrcMAC [6]byte
	Type   [2]byte
	Data   []byte
}

type ArpPacket struct {
	HardwareType          [2]byte
	ProtocolType          [2]byte
	HardwareLen           [1]byte
	ProtocolLen           [1]byte
	Operation             [2]byte
	SenderHWAddress       [6]byte
	SenderProtocolAddress [4]byte
	TargetHWAddress       [6]byte
	TargetProtocolAddress [4]byte
}

func NewArpRequest(srcAddr, dstAddr [6]byte, srcIP, dstIP [4]byte) *ArpPacket {
	return &ArpPacket{
		HardwareType:          [2]byte{0x00, 0x01},
		ProtocolType:          [2]byte{0x08, 0x00},
		HardwareLen:           [1]byte{0x06},
		ProtocolLen:           [1]byte{0x04},
		Operation:             [2]byte{0x00, 0x01},
		SenderHWAddress:       srcAddr,
		SenderProtocolAddress: srcIP,
		TargetHWAddress:       dstAddr,
		TargetProtocolAddress: dstIP,
	}
}

func (a *ArpPacket) ToPacket() []byte {
	ret := make([]byte, 0)
	buf := bytes.NewBuffer(ret)
	buf.Write(a.HardwareType[:])
	buf.Write(a.ProtocolType[:])
	buf.Write(a.HardwareLen[:])
	buf.Write(a.ProtocolLen[:])
	buf.Write(a.Operation[:])
	buf.Write(a.SenderHWAddress[:])
	buf.Write(a.SenderProtocolAddress[:])
	buf.Write(a.TargetHWAddress[:])
	buf.Write(a.TargetProtocolAddress[:])
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
			DstMAC: [6]byte{},
			SrcMAC: [6]byte{},
			Type:   [2]byte{},
			Data:   []byte{},
		},
		FCS: [4]byte{0, 0, 0, 0},
	}
}

func (e *Ethernet) Send() error {
	return nil
}
