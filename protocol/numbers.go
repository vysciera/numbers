package protocol

import (
	"encoding/binary"
	"errors"
)

const NumberGroupSize = 10

type NumberGroup [5]uint16

func EncodeNumberGroup(group NumberGroup) []byte {
	buf := make([]byte, NumberGroupSize)

	for i, n := range group {
		offset := i * 2

		binary.BigEndian.PutUint16(
			buf[offset:offset + 2],
			n,
		)
	}

	return buf
}

func DecodeNumberGroup(data []byte) (NumberGroup, error) {
	var group NumberGroup

	if len(data) != NumberGroupSize {
		return group, errors.New("invalid number group size")
	}

	for i := range group {
		offset := i * 2

		group[i] = binary.BigEndian.Uint16(data[offset : offset + 2])
	}

	return group, nil
}

