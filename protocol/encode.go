package protocol

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
)

func Encode(p Packet) ([]byte, error) {
	if err := validatePacket(p); err != nil {
		return nil, err
	}

	payloadLength := len(p.Payload)
	packetLength := HeaderSize + payloadLength + CRCSize

	buf := make([]byte, packetLength)

	// 0..3 - magic
	copy(buf[0:4], Magic[:])

	// 4 - version
	buf[4] = p.Version

	// 5 - type
	buf[5] = byte(p.Type)

	// 6 - flags
	buf[6] = byte(p.Flags)

	// 7 - header length
	buf[7] = HeaderSize

	// 8..9 - station ID
	binary.BigEndian.PutUint16(
		buf[8:10],
		p.StationID,
	)

	// 10..13 - global sequence
	binary.BigEndian.PutUint32(
		buf[10:14],
		p.Sequence,
	)

	// 14..17 - transmission ID
	binary.BigEndian.PutUint32(
		buf[14:18],
		p.TransmissionID,
	)

	// 18..19 - transmission index
	binary.BigEndian.PutUint16(
		buf[18:20],
		p.Index,
	)

	// 22..29 - timestamp
	binary.BigEndian.PutUint16(
		buf[20:22],
		p.Count,
	)

	// 30..31 - payload length
	binary.BigEndian.PutUint16(
		buf[30:32],
		uint16(payloadLength),
	)

	// 32..N - payload
	copy(
		buf[HeaderSize:HeaderSize+payloadLength],
		p.Payload,
	)

	// final four bytes - CRC32
	checksumOffset := HeaderSize + payloadLength
	checksum := crc32.ChecksumIEEE(buf[:checksumOffset])
	binary.BigEndian.PutUint32(
		buf[checksumOffset:checksumOffset+CRCSize],
		checksum,
	)

	return buf, nil
}

func validatePacket(p Packet) error {
	if p.Version != Version1 {
		return errors.New("unsupported NSP version")
	}

	if !p.Type.Valid() {
		return errors.New("invalid packet type")
	}

	if p.StationID == 0 {
		return errors.New("station ID 0 is reserved")
	}

	if len(p.Payload) > MaxPayload {
		return errors.New("payload exceeds NSP/1 maximum")
	}

	switch p.Type {
	case TypeBeacon:
		if p.TransmissionID != 0 {
			return errors.New("beacon cannot belong to a transmission")
		}

		if p.Index != 0 || p.Count != 0 {
			return errors.New("beacon index/count must be zero")
		}

		if len(p.Payload) != 0 {
			return errors.New("beacon cannot contain payload")
		}
	
	case TypePreamble:
		if p.TransmissionID == 0 {
			return errors.New("preamble requires transmission ID")
		}

		if p.Index != 0 {
			return errors.New("preamble index must be zero")
		}

		if len(p.Payload) != 0 {
			return errors.New("preamble cannot contain payload")
		}

	case TypeMessage:
		if p.TransmissionID == 0 {
			return errors.New("message requires transmission ID")
		}

		if p.Count == 0 {
			return errors.New("message count cannot be zero")
		}

		if p.Index == 0 || p.Index > p.Count {
			return errors.New("invalid message index")
		}

	case TypeRepeat:
		if p.TransmissionID == 0 {
			return errors.New("repeat requires transmission ID")
		}

		if p.Index != 0 {
			return errors.New("repeat index must be zero")
		}

		if len(p.Payload) != 0 {
			return errors.New("repeat cannot contain payload")
		}

	case TypeTerminator:
		if p.TransmissionID == 0 {
			return errors.New("terminator requires transmission ID")
		}

		if p.Index != 0 {
			return errors.New("terminator index must be zero")
		}

		if len(p.Payload) != 0 {
			return errors.New("terminator cannot contain payload")
		}
	}

	return nil
}
