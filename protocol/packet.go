package protocol

type Packet struct {
	Version	uint8
	Type	Type
	Flags	Flags

	StationID	uint16
	Sequence	uint32

	TransmissionID	uint32
	Index			uint16
	Count			uint16

	Timestamp	uint64
	Payload	[]byte
}
