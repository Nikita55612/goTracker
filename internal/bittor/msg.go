package bittor

import "encoding/binary"

type Tag int

const (
	Choke Tag = iota
	Unchoke
	Interested
	NotInterested
	Have
	Bitfield
	Request
	Piece
	Cancel
	Unknown
)

func TagFromByte(b byte) Tag {
	switch {
	case b <= 8:
		return Tag(b)
	default:
		return Tag(9)
	}
}

type Message struct {
	Tag     Tag
	Payload []byte
}

func MessageFromBytes(b []byte) *Message {
	if len(b) <= 1 {
		return nil
	}
	length := binary.BigEndian.Uint32(b[:4])
	return &Message{
		Tag:     TagFromByte(b[4]),
		Payload: b[5:length],
	}
}

func (m *Message) Serialize() []byte {
	if m == nil {
		return make([]byte, 4)
	}
	length := uint32(len(m.Payload) + 1)
	buf := make([]byte, 4+length)
	binary.BigEndian.PutUint32(buf[0:4], length)
	buf[4] = byte(m.Tag)
	copy(buf[5:], m.Payload)
	return buf
}
