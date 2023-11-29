package main

import (
	"net"
	"os"
	"time"

	"atomicgo.dev/keyboard/keys"
	"github.com/charmbracelet/log"

	"github.com/spf13/cobra"

	"github.com/stelmanjones/wrc"
	"github.com/stelmanjones/wrc/cmd/tui/input"
	"github.com/stelmanjones/wrc/cmd/webserver/routes"
	"github.com/valyala/fasthttp"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 20 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var rootCmd = &cobra.Command{
	Use:   "wrc",
	Short: "A telemetry server and library made for EAS WRC.",
	Run: func(cmd *cobra.Command, _args []string) {
		run()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Error(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}

func run() {
	ch := make(chan wrc.Packet)
	in := make(chan keys.Key)

	shutdown := func() {
		close(ch)
		close(in)
		log.Info("Bye!")
		os.Exit(0)
	}

	conn, err := net.ListenPacket("udp4", ":6969")
	if err != nil {
		log.Error(err)
	}
	defer conn.Close()

	app := routes.RegisterRoutes()
	go routes.Client.Run(conn)
	go fasthttp.ListenAndServe(":9999", app.Handler)
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

				err := routes.Client.Push(&packet)
				if err != nil {
					log.Error(err)
				}

			}
		default:
			continue
		}
	}
}

