package protocol

import "errors"

func Validate(p Packet) error {
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

		if p.Count == 0 {
			return errors.New("preamble count cannot be zero")
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

		if len(p.Payload) != NumberGroupSize {
			return errors.New("message requires one number-group payload")
		}

	case TypeRepeat:
		if p.TransmissionID == 0 {
			return errors.New("repeat requires transmission ID")
		}

		if p.Index != 0 {
			return errors.New("repeat index must be zero")
		}

		if p.Count == 0 {
			return errors.New("repeat count cannot be zero")
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

		if p.Count == 0 {
			return errors.New("terminator count cannot be zero")
		}

		if len(p.Payload) != 0 {
			return errors.New("terminator cannot contain payload")
		}
	}

	return nil
}
