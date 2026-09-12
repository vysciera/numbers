package protocol

func NewBeacon(stationID uint16, sequence uint32, timestamp uint64) (Packet, error) {
	p := Packet{
		Version:	Version1,
		Type:		TypeBeacon,
		StationID:	stationID,
		Sequence:	sequence,
		Timestamp:	timestamp,
	}

	if err := Validate(p); err != nil {
		return Packet{}, err
	}

	return p, nil
}

func NewPreamble(stationID, count uint16, transmissionID, sequence uint32, timestamp uint64) (Packet, error) {
	p := Packet{
		Version:	Version1,
		Type:		TypePreamble,
		StationID:	stationID,
		Sequence:	sequence,
		TransmissionID:	transmissionID,
		Count:	count,
		Timestamp:	timestamp,
	}

	if err := Validate(p); err != nil {
		return Packet{}, err
	}

	return p, nil
}

func NewMessage(stationID, count, index uint16, transmissionID, sequence uint32, timestamp uint64, group NumberGroup, repeated bool) (Packet, error) {
	var flags Flags

	if repeated {
		flags |= FlagRepeated
	}

	p := Packet{
		Version:	Version1,
		Type:	TypeMessage,
		Flags: flags,
		StationID: stationID,
		Sequence: sequence,
		TransmissionID:	transmissionID,
		Index:	index,
		Count:	count,
		Timestamp:	timestamp,
		Payload:	EncodeNumberGroup(group),
	}

	if err := Validate(p); err != nil {
		return Packet{}, err
	}

	return p, nil
}

func NewRepeat(stationID, count uint16, transmissionID, sequence uint32, timestamp uint64) (Packet, error) {
	p := Packet{
		Version:	Version1,
		Type:	TypeRepeat,
		StationID:	stationID,
		Sequence:	sequence,
		TransmissionID:	transmissionID,
		Count:	count,
		Timestamp:	timestamp,
	}

	if err := Validate(p); err != nil {
		return Packet{}, err
	}

	return p, nil
}

func NewTerminator(stationID, count uint16, transmissionID, sequence uint32, timestamp uint64) (Packet, error) {
	p := Packet{
		Version:	Version1,
		Type:	TypeTerminator,
		StationID:	stationID,
		Sequence:	sequence,
		TransmissionID:	transmissionID,
		Count:	count,
		Timestamp:	timestamp,
	}

	if err := Validate(p); err != nil {
		return Packet{}, err
	}

	return p, nil
}
