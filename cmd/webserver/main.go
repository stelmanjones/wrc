package main

import (
	"os"

	"atomicgo.dev/keyboard/keys"
	"github.com/charmbracelet/log"

	"github.com/spf13/cobra"

	"github.com/stelmanjones/wrc"
	"github.com/stelmanjones/wrc/cmd/tui/input"
	"github.com/stelmanjones/wrc/cmd/webserver/routes"
	"github.com/valyala/fasthttp"
)

var (
	rootCmd = &cobra.Command{
		Use:   "wrc",
		Short: "A telemetry server and library made for EAS WRC.",
		Run: func(_ *cobra.Command, _ []string) {
			run()
		},
	}

	debugFlag = false
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "d", false, "Set debug flag")
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
	if debugFlag {
		log.SetLevel(log.DebugLevel)
	}
	ch := make(chan wrc.Packet)
	in := make(chan keys.Key)
	routes.Client.Debug = debugFlag

	shutdown := func() {
		close(ch)
		close(in)
		log.Fatal("Bye! 👋")
		os.Exit(0)
	}

	app := routes.RegisterRoutes()
	go routes.Client.Run()
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
				} else {
					log.Debug("Pushed new packet.")
				}

			}
		default:
			continue
		}
	}
}
