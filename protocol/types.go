package protocol

const (
	Version1 uint8 = 1

	HeaderSize  = 32
	CRCSize 	= 4
	MaxPacket 	= 512
	MaxPayload	= MaxPacket - HeaderSize - CRCSize
)

var Magic = [4]byte{'N', 'S', 'P', '1'}

type Type uint8

const(
	TypeBeacon		Type = 0x01
	TypePreamble	Type = 0x02
	TypeMessage		Type = 0x03
	TypeRepeat		Type = 0x04
	TypeTerminator	Type = 0x05
)

func (t Type) Valid() bool {
	switch t {
		case TypeBeacon,
			TypePreamble,
			TypeMessage,
			TypeRepeat,
			TypeTerminator:
			return true
		default:
			return false
	}
}

type Flags uint8

const (
	FlagRepeated Flags = 1 << 0
)

func(f Flags) Has(flag Flags) bool {
	return f&flag != x01
}
