package station

import (
	"errors"
	"time"

	"numbers/protocol"
	"numbers/transmission"
	"numbers/transport"
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

func (s *Station) emit(sender transport.Sender, packet protocol.Packet) (protocol.Packet, error) {
	raw, err := protocol.Encode(packet)
	if err != nil {
		return protocol.Packet{}, err
	}

	if err := sender.Send(raw); err != nil {
		return protocol.Packet{}, err
	}

	s.commitSequence(packet.Sequence)

	return packet, nil
}

func (s *Station) EmitBeacon(sender transport.Sender) (protocol.Packet, error) {
	sequence := s.candidateSequence()

	packet, err := protocol.NewBeacon(s.ID, sequence, now())
	if err != nil {
		return protocol.Packet{}, err
	}

	return s.emit(sender, packet)
}

func (s *Station) EmitFrame(sender transport.Sender, frame transmission.Frame) (protocol.Packet, error) {
	sequence := s.candidateSequence()
	timestamp := now()

	var (
		packet protocol.Packet
		err	   error
	)

	switch frame.Type {
	case protocol.TypePreamble:
		packet, err = protocol.NewPreamble(
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

		packet, err = protocol.NewMessage(
			s.ID,
			frame.Count,
			frame.Index,
			frame.TransmissionID,
			sequence,
			timestamp,
			*frame.Group,
			frame.Flags.Has(protocol.FlagRepeated),
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
		return protocol.Packet{},
		errors.New("unsupported tnansmission frame")
	}

	if err != nil {
		return protocol.Packet{}, err
	}

	return s.emit(sender, packet)
}

