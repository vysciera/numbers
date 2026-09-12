package station

import (
	"errors"
	"time"

	"numbers/protocol"
	"numbers/transmission"
)

type Station struct {
	ID uint16

	sequence uint32
}

func New(id uint16) (*Station, error) {
	if id == 0 {
		return nil, errors.New("station ID 0 is reserved")
	}

	return &Station{
		ID: id,
	}, nil
}

func (s *Station) candidateSequence() uint32 {
	return s.sequence + 1
}

func (s *Station) commitSequence(sequence uint32) {
	s.sequence = sequence
}

func now() uint64 {
	return uint64(time.Now().UnixMilli())
}

func (s *Station) Beacon() (protocol.Packet, error) {
	sequence := s.candidateSequence()

	packet, err := protocol.NewBeacon(s.ID, sequence, now())
	if err != nil {
		return protocol.Packet{}, err
	}

	s.commitSequence(sequence)
	
	return packet, nil
}

func (s *Station) Packet(frame transmission.Frame) (protocol.Packet, error) {
	sequence := s.candidateSequence()
	timestamp := now()

	var (
		packet protocol.Packet
		err	   error
	)

	switch frame.Type {
	case protocol.TypePreamble:
		packet, err =  protocol.NewPreamble(
			s.ID,
			frame.Count,
			frame.TransmissionID,
			sequence,
			timestamp,
		)

	case protocol.TypeMessage:
		if frame.Group == nil {
			return protocol.Packet{}, errors.New("message frame has no number group")
		}

		repeated := frame.Flags.Has(protocol.FlagRepeated)

		packet, err = protocol.NewMessage(
			s.ID,
			frame.Count,
			frame.Index,
			frame.TransmissionID,
			sequence,
			timestamp,
			*frame.Group,
			repeated,
		)

	case protocol.TypeRepeat:
		packet, err = protocol.NewRepeat(
			s.ID,
			frame.Count,
			frame.TransmissionID,
			sequence,
			timestamp,
		)

	case protocol.TypeTerminator:
		packet, err = protocol.NewTerminator(
			s.ID,
			frame.Count,
			frame.TransmissionID,
			sequence,
			timestamp,
		)

	default:
		return protocol.Packet{}, errors.New("unsupported transmission frame")
	}

	if err != nil {
		return protocol.Packet{}, err
	}

	s.commitSequence(sequence)

	return packet, nil
}
