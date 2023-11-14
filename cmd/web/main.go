package main

import (
	"bufio"
	"embed"
	"fmt"
	"html/template"
	"net"
	"os"
	"time"

	"atomicgo.dev/keyboard/keys"
	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/valyala/fasthttp"

	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/spf13/pflag"
	"github.com/stelmanjones/wrc"
	"github.com/stelmanjones/wrc/cmd/tui/input"
	"github.com/stelmanjones/wrc/cmd/web/routes"
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

	app.Use(routes.CorsHandler())
	app.Use("/assets", routes.EmbedFSHandler(&embedDirAssets))
	app.Get("/metrics", monitor.New(monitor.Config{Title: "WRC Web Telemetry Metrics"}))
	app.Get("/", routes.RootHandler())


	tplGroup := app.Group("templates")
	tplGroup.Get("timer", TimerHandler())



	apiGroup := app.Group("api")
	apiGroup.Get("packet", func(c *fiber.Ctx) error {
		c.JSON(WrcPacket.Data)
		return nil
	})

	sseGroup := app.Group("sse")
	sseGroup.Get("/sse",SSEHandler())

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
	go udp.ListenForPacket(conn, ch)
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
				if packet.PacketUID != WrcPacket.Data.PacketUID {
				WrcPacket.Data = packet
				}
				WrcPacket.Mu.Unlock()
			}
		default:
			continue
		}
	}
}

func TimerHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		WrcPacket.Mu.RLock()
		defer WrcPacket.Mu.RUnlock()
		//	WrcPacket.Data.PacketUID ++
		// WrcPacket.Data.GameTotalTime ++
		// WrcPacket.Data.StageCurrentTime ++
		// WrcPacket.Data.StageLength = 23.46
		// WrcPacket.Data.StageCurrentDistance += 0.01
		c.Render("partials/timer", fiber.Map{
		})
		return nil
	}
}

func SSEHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			log.Info("SSE Client connected!")
			for {
			
				fmt.Fprintf(w, "event: currentTime\ndata: %s\n\n", WrcPacket.Data.CurrentStageTime())

				err := w.Flush()
				if err != nil {
					// Refreshing page in web browser will establish a new
					// SSE connection, but only (the last) one is alive, so
					// dead connections must be closed here.
					log.Errorf("Error while flushing: %v. Closing http connection.\n", err)

					break
				}
				fmt.Fprintf(w, "event: inGameTime\ndata: %s\n\n", WrcPacket.Data.InGameTime())

				err = w.Flush()
				if err != nil {
					// Refreshing page in web browser will establish a new
					// SSE connection, but only (the last) one is alive, so
					// dead connections must be closed here.
					log.Errorf("Error while flushing: %v. Closing http connection.\n", err)

					break
				}
				time.Sleep(100 * time.Millisecond)
			}
		}))
		return nil
	}

}