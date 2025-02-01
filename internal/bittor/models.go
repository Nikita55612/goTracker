package bittor

import (
	"encoding/binary"
	"fmt"
	"gotracker/internal/tracker"
	"net"
)

type Torrent struct {
	Announce string
	Info     tracker.BencodeInfo
	InfoHash tracker.InfoHash
	Peers    Peers
	Interval int
}

type Peer struct {
	Ip   net.IP
	Port uint16
}

func (p *Peer) String() string {
	return fmt.Sprintf(
		"%s:%d", p.Ip, p.Port,
	)
}

type Peers []Peer

func PeersFromBytes(b []byte) Peers {
	numPeers := len(b) / 6
	peers := make(Peers, numPeers)
	for i := 0; i < numPeers; i++ {
		offset := i * 6
		peers[i].Ip = net.IP(b[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16(b[offset+4 : offset+6])
	}
	return peers
}
