package main

import (
	"fmt"
	"log"
	"time"

	"numbers/protocol"
	"numbers/station"
	"numbers/transmission"
	"numbers/transport"
)

const address = "127.0.0.1:4040"

func main() {
	udp, err := transport.NewUDP(address)
	if err != nil {
		log.Fatal(err)
	}
	defer udp.Close()

	const statePath = "var/station.json"

	s, err := station.Open(statePath, 17)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf(
		"station %04d transmitting to %s\n\n",
		s.ID(), address,
	)

	packet, err := s.EmitBeacon(udp)
	if err != nil {
		log.Fatal(err)
	}

	printPacket(packet)

	time.Sleep(2 * time.Second)

	transmissionID, err := s.NextTransmissionID()
	if err != nil {
		log.Fatal(err)
	}

	tx := transmission.Transmission{
		ID: transmissionID,

		Groups: []protocol.NumberGroup{
			{418, 992, 117, 4, 731},
			{882, 119, 500, 771, 229},
			{401, 401, 793, 8, 144},
		},
	}

	frames, err := tx.Frames()
	if err != nil {
		log.Fatal(err)
	}

	for _, frame := range frames {
		packet, err := s.EmitFrame(udp, frame)
		if err != nil {
			log.Fatal(err)
		}

		printPacket(packet)

		time.Sleep(2 * time.Second)
	}

	packet, err = s.EmitBeacon(udp)
	if err != nil {
		log.Fatal(err)
	}

	printPacket(packet)
}

func printPacket(p protocol.Packet) {
	fmt.Printf(
		"%08d %-10s station=%04d",
		p.Sequence, p.Type, p.StationID,
	)

	if p.TransmissionID != 0 {
		fmt.Printf(" tx=%d", p.TransmissionID)
	}

	if p.Type == protocol.TypeMessage {
		group, err := protocol.DecodeNumberGroup(p.Payload)
		if err != nil {
			fmt.Printf(" INVALID PAYLOAD\n")
			return
		}

		fmt.Printf(
			" group=%d/%d %03d %03d %03d %03d %03d",
			p.Index, p.Count,
			group[0], group[1], group[2], group[3], group[4],
		)

		if p.Flags.Has(protocol.FlagRepeated) {
			fmt.Print(" REPEATED")
		}
	}

	fmt.Println()
}
