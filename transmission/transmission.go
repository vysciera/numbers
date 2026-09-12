package transmission

import (
	"errors"
	
	"numbers/protocol"
)

type Transmission struct {
	ID	uint32
	Groups []protocol.NumberGroup
}

type Frame struct {
	Type	protocol.Type
	Flags	protocol.Flags

	TransmissionID	uint32

	Index	uint16
	Count	uint16

	Group *protocol.NumberGroup
}

func (t Transmission) Frames() ([]Frame, error) {
	if t.ID == 0 {
		return nil, errors.New("transmission ID 0 is reserved")
	}

	if len(t.Groups) == 0 {
		return nil, errors.New("transmission reserves at least one group")
	}

	if len(t.Groups) > 65535 {
		return nil, errors.New("too many groups")
	}

	count := uint16(len(t.Groups))
	frames := make([]Frame, 0, 2 * len(t.Groups) + 3)

	// Preamble
	frames = append(frames, Frame{
		Type:	protocol.TypePreamble,
		TransmissionID:	t.ID,
		Count:	count,
	})

	// Original pass
	for i := range t.Groups {
		group := t.Groups[i]
		frames = append(frames, Frame{
			Type:	protocol.TypeMessage,
			TransmissionID: t.ID,
			Index:	uint16(i + 1),
			Count:	count,
			Group:	&group,
		})
	}

	// REPEAT marker
	frames = append(frames, Frame{
		Type:	protocol.TypeRepeat,
		TransmissionID: t.ID,
		Count:	count,
	})

	// Repeated pass 
	for i := range t.Groups {
		group := t.Groups[i]
		frames = append(frames, Frame{
			Type:	protocol.TypeMessage,
			Flags:	protocol.FlagRepeated,
			TransmissionID:	t.ID,
			Index:	uint16(i + 1),
			Count:	count,
			Group:	&group,
		})
	}

	// Terminator
	frames = append(frames, Frame{
		Type:	protocol.TypeTerminator,
		TransmissionID: t.ID,
		Count:	count,
	})

	return frames, nil
}
