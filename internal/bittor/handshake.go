package bittor

import "gotracker/internal/tracker"

type Handshake [68]byte

func NewHandshake(
	infoHash tracker.InfoHash,
	peerId tracker.PeerId,
) Handshake {
	const pstr = "BitTorrent protocol"
	handshake := [68]byte{}
	handshake[0] = byte(len(pstr))
	copy(handshake[1:], pstr)
	copy(handshake[1+len(pstr)+8:], infoHash[:])
	copy(handshake[1+len(pstr)+8+20:], peerId[:])
	return handshake
}
