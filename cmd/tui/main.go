package main

import (
	"net"
	"os"

	"atomicgo.dev/keyboard/keys"
	"github.com/charmbracelet/log"

	flag "github.com/spf13/pflag"
	"github.com/stelmanjones/wrc"
	"github.com/stelmanjones/wrc/cmd/tui/input"
	"github.com/stelmanjones/wrc/server/udp"
)

var (
	udpAddress  string
	apiAddress  string
	refreshRate int
)

func main() {
	flag.StringVarP(&udpAddress, "udp", "u", ":6969", "Set UDP connection address.")
	flag.StringVar(&apiAddress, "api", ":9999", "Set API server address.")
	flag.IntVarP(&refreshRate, "refresh-rate", "r", 10, "Set udp refresh rate.")
	flag.Parse()

	log.SetLevel(log.DebugLevel)

	in := make(chan keys.Key)
	ch := make(chan wrc.Packet)
	shutdown := func() {
		close(in)
		close(ch)
		os.Exit(0)
	}

	conn, err := net.ListenPacket("udp4", udpAddress)
	if err != nil {
		log.Error(err)
	}

	defer conn.Close()
	log.Debug("Starting server!", "address", udpAddress)
	go udp.ListenForPacket(conn, ch, refreshRate)
	go input.ListenForInput(in)

	var packet wrc.Packet
	for {
		select {
		case key := <-in:
			{
				switch key.Code {
				case keys.RuneKey:
					{
						switch key.String() {
						case "q":
							{
								shutdown()
							}
						default:
							{
								continue
							}
						}
					}
				case keys.CtrlC, keys.Escape:
					{
						shutdown()
					}
				default:
					{
						continue
					}
				}
			}
		case packet = <-ch:
			{
				log.Infof("Packet:", "Stage Length", packet.StageLength, "Current Time", packet.StageCurrentTime)
			}
		default:
			continue
		}
	}
}
