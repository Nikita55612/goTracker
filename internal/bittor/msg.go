package bittor

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
	return &Message{
		Tag:     TagFromByte(b[0]),
		Payload: b[1:],
	}
}
