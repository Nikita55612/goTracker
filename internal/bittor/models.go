package bittor

import (
	"encoding/binary"
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

type Peers []Peer

func PeersFromBytes(b []byte) Peers {
	const peerSize = 6
	numPeers := len(b) / peerSize
	peers := make(Peers, numPeers)
	for i := 0; i < numPeers; i++ {
		offset := i * peerSize
		peers[i].Ip = net.IP(b[offset : offset+4])
		peers[i].Port = binary.BigEndian.Uint16(b[offset+4 : offset+6])
	}
	return peers
}
