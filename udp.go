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

//Listen is a function that listens for incoming packets on a given connection.
//It reads the packets into a buffer, decodes them into a Packet struct, and sends them on a channel.
//
//	go Listen(conn,ch)
// 	ch: The channel to send the decoded packets on.
//
func Listen(conn net.PacketConn, ch chan Packet) error {
	buf := make([]byte, binary.Size(Packet{}))
	for {
		_, _, err := conn.ReadFrom(buf)
		if err != nil {
			logger.Error(err)
			return ErrUDPData
		}
		var packet Packet
		err = binary.Read(bytes.NewReader(buf), binary.LittleEndian, &packet)
		if err != nil {
			logger.Error(err)
			return ErrUDPData
		}

		ch <- packet

	}
	
}
