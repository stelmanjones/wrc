package udp

import (
	"bytes"
	"encoding/binary"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/stelmanjones/wrc"
)

var logger log.Logger = *log.New(os.Stdout).WithPrefix("UDP")

func ListenForPacket(conn net.PacketConn, ch chan wrc.Packet, refreshRate int) {
	buf := make([]byte, binary.Size(wrc.Packet{}))
	delay := 1000 / refreshRate
	ticker := time.NewTicker(time.Duration(delay * int(time.Millisecond)))
	defer ticker.Stop()
	for {
		_, _, err := conn.ReadFrom(buf)
		if err != nil {
			logger.Error(err)
		}
		var packet wrc.Packet
		err = binary.Read(bytes.NewReader(buf), binary.LittleEndian, &packet)
		if err != nil {
			logger.Error(err)
		}
		ch <- packet

	}
}
