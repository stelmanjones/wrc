package main

import (
	"embed"
	"html/template"
	"net"
	"os"

	"atomicgo.dev/keyboard/keys"
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/spf13/pflag"
	"github.com/stelmanjones/wrc"
	"github.com/stelmanjones/wrc/cmd/tui/input"
	"github.com/stelmanjones/wrc/cmd/web/endpoints"
	"github.com/stelmanjones/wrc/server/api"
	"github.com/stelmanjones/wrc/server/udp"
)

var (
	//go:embed assets/*
	embedDirAssets embed.FS
	udpAddress     string
	apiAddress     string
	configPath     string
	refreshRate    int

	WrcPacket *wrc.ThreadSafePacket = wrc.NewThreadSafePacket()
)

func main() {
	debugEnv := os.Getenv("ENV")

	debugEnabled := func() bool {
		switch debugEnv {
		case "debug":
			return true
		default:
			return false
		}
	}()
	if debugEnabled {

		log.SetLevel(log.DebugLevel)
		udpAddress = ":6969"
	}

	pflag.StringVarP(&udpAddress, "udp", "u", ":6969", "Set UDP connection address.")
	pflag.StringVar(&apiAddress, "api", ":9999", "Set API server address.")
	pflag.IntVarP(&refreshRate, "refresh-rate", "r", 10, "Set udp refresh rate.")
	pflag.Parse()

	engine := html.New("cmd/web/templates", ".tmpl")
	engine.AddFunc(
		// add unescape function
		"unescape", func(s string) template.HTML {
			return template.HTML(s)
		},
	)
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Use(endpoints.CorsHandler())
	app.Use("/assets", endpoints.EmbedFSHandler(&embedDirAssets))
	app.Get("/metrics", monitor.New(monitor.Config{Title: "WRC Web Telemetry Metrics"}))
	app.Get("/", endpoints.RootHandler())

	tplGroup := app.Group("templates")
	tplGroup.Get("timer", TimerHandler())
	
	apiGroup := app.Group("api")
	apiGroup.Get("packet", func(c *fiber.Ctx) error {
		c.JSON(WrcPacket.Packet)
		return nil

	})

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

	go api.RunHttpServer(app, ":9999", WrcPacket)

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
				WrcPacket.Mu.Lock()
				WrcPacket.Packet = packet
				WrcPacket.Mu.Unlock()
			}
		default:
			continue
		}
	}
}

func TimerHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Render("partials/timer", fiber.Map{
			"Current": WrcPacket.Packet.CurrentStageTime(),
			"Total": WrcPacket.Packet.InGameTime(),
		})
		return nil
	}
}
