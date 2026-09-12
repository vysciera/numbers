package main

import (
	"fmt"
	"log"

	"numbers/protocol"
	"numbers/transmission"
	"numbers/station"
)

func main() {
	s, err := station.New(17)
	if err != nil {
		log.Fatal(err)
	}

	// Station init, idle
	beacon, err := s.Beacon()
	if err != nil {
		log.Fatal(err)
	}

	printPacket(beacon)

	tx := transmission.Transmission{
		ID: 481,

		Groups: []protocol.NumberGroup{
			{ 418, 992, 117, 4, 731 },
			{ 882, 119, 500, 771, 229 },
			{ 401, 401, 793, 8, 114 },
		},
	}

	frames, err := tx.Frames()
	if err != nil {
		log.Fatal(err)
	}

	for _, frame := range frames {
		packet, err := s.Packet(frame)
		if err != nil {
			log.Fatal(err)
		}

		printPacket(packet)
	}

	// Transmission complete: idle
	beacon, err = s.Beacon()
	if err != nil {
		log.Fatal(err)
	}

	printPacket(beacon)
}

func printPacket(p protocol.Packet) {
	fmt.Printf(
		"%08d %-10s station=%04d",
		p.Sequence,
		p.Type,
		p.StationID,
	)

	if p.TransmissionID != 0 {
		fmt.Printf(
			" tx=%d",
			p.TransmissionID,
		)
	}

	if p.Type == protocol.TypeMessage {
		group, err := protocol.DecodeNumberGroup(p.Payload)
		if err != nil {
			fmt.Printf(" INVALID PAYLOAD\n")
			return
		}

		fmt.Printf(
			" group=%d/%d %03d %03d %03d %03d %03d",
			p.Index,
			p.Count,
			group[0],
			group[1],
			group[2],
			group[3],
			group[4],
		)

		if p.Flags.Has(protocol.FlagRepeated) {
			fmt.Printf(" REPEATED")
		}
	}

	fmt.Println()
}
