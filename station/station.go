package station 

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"numbers/protocol"
	"numbers/transmission"
	"numbers/transport"
)

type Station struct {
	mu sync.Mutex

	path	string
	state	State
}

func Open(path string, stationID uint16) (*Station, error) {
	state, err := loadState(path)

	switch {
		case err == nil:
			if state.StationID != stationID {
				return nil, fmt.Errorf(
					"state belongs to station %d, not station %d",
					state.StationID, stationID,
				)
			}

			state.Starts++
			state.LastStartedAt=time.Now().UTC()

		case errors.Is(err, os.ErrNotExist):
			state, err = newState(stationID)
			if err != nil {
				return nil, err
			}

		default:
			return nil, err
	}

	if err := saveState(path, state); err != nil {
		return nil, err
	}

	return &Station{
		path:	path,
		state:	state,
	}, nil
}

func (s *Station) ID() uint16 {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.state.StationID
}

func (s *Station) Epoch() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.state.Epoch
}

func (s *Station) Snapshot() State {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.state
}

func (s *Station) candidateSequence() uint32 {
	return s.state.Sequence + 1
}

func (s *Station) reserveSequence(sequence uint32) error {
	old := s.state.Sequence
	s.state.Sequence = sequence

	if err := saveState(s.path, s.state); err != nil {
		s.state.Sequence = old
		return err
	}

	return nil
}

func (s *Station) NextTransmissionID() (uint32, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	next := s.state.Transmission + 1

	if next == 0 {
		next = 1
	}

	old := s.state.Transmission
	s.state.Transmission = next

	if err := saveState(s.path, s.state); err != nil {
		s.state.Transmission = old
		return 0, err
	}

	return next, nil
}

func (s *Station) emitLocked(sender transport.Sender, packet protocol.Packet) (protocol.Packet, error) {
	raw, err := protocol.Encode(packet)
	if err != nil {
		return protocol.Packet{}, err
	}

	if err := s.reserveSequence(packet.Sequence); err != nil {
		return protocol.Packet{}, fmt.Errorf(
			"reserve sequence %d: %w",
			packet.Sequence, err,
		)
	}

	if err := sender.Send(raw); err != nil {
		return packet, fmt.Errorf(
			"emit sequence %d: %w",
			packet.Sequence, err,
		)
	}

	return packet, nil
}

func (s *Station) EmitBeacon(sender transport.Sender) (protocol.Packet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sequence := s.candidateSequence()
	packet, err := protocol.NewBeacon(
		s.state.StationID, sequence,
		uint64(time.Now().UnixMilli()),
	)
	
	if err != nil {
		return protocol.Packet{}, err
	}

	return s.emitLocked(sender, packet)
}

func (s *Station) EmitFrame(sender transport.Sender, frame transmission.Frame) (protocol.Packet, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	sequence := s.candidateSequence()
	timestamp := uint64(time.Now().UnixMilli())

	var (
		packet	protocol.Packet
		err		error
	)

	switch frame.Type {
		case protocol.TypePreamble:
			packet, err = protocol.NewPreamble(
				s.state.StationID,
				frame.Count,
				frame.TransmissionID,
				sequence,
				timestamp,
			)

		case protocol.TypeMessage:
			if frame.Group == nil {
				return protocol.Packet{},
				errors.New("message frame has no number group")
			}

			packet, err = protocol.NewMessage(
				s.state.StationID,
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
				s.state.StationID,
				frame.Count,
				frame.TransmissionID,
				sequence,
				timestamp,
			)

		case protocol.TypeTerminator:
			packet, err = protocol.NewTerminator(
				s.state.StationID,
				frame.Count,
				frame.TransmissionID,
				sequence,
				timestamp,
			)

		default:
			return protocol.Packet{},
			errors.New("unsupported transmission frame")
	}

	if err != nil {
		return protocol.Packet{}, err
	}

	return s.emitLocked(sender, packet)
}
