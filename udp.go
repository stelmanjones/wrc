package wrc

import (
	"bytes"
	"encoding/binary"
	"net"
	"os"

	"github.com/charmbracelet/log"
)

var logger = log.NewWithOptions(os.Stderr, log.Options{
	Prefix: "UDP",
})

func ListenForPacket(conn net.PacketConn, ch chan Packet) {
	buf := make([]byte, binary.Size(Packet{}))
	for {
		_, _, err := conn.ReadFrom(buf)
		if err != nil {
			logger.Error(err)
		}
		var packet Packet
		err = binary.Read(bytes.NewReader(buf), binary.LittleEndian, &packet)
		if err != nil {
			logger.Error(err)
		}

		ch <- packet

	}
}
