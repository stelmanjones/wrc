package wrc

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"

	"github.com/charmbracelet/log"
)

func ListenForPacket(conn net.PacketConn, ch chan Packet) {
	buf := make([]byte, binary.Size(Packet{}))
	delay := 1000 / 30
	ticker := time.NewTicker(time.Duration(delay * int(time.Millisecond)))
	defer ticker.Stop()
	for {
		_, _, err := conn.ReadFrom(buf)
		if err != nil {
			log.WithPrefix("UDP").Error(err)
		}
		var packet Packet
		err = binary.Read(bytes.NewReader(buf), binary.LittleEndian, &packet)
		if err != nil {
			log.WithPrefix("UDP").Error(err)
		}
		packet.PacketUID++
		ch <- packet

	}
}
