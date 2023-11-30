package routes

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/fasthttp/router"
	"github.com/fasthttp/websocket"
	"github.com/stelmanjones/wrc"
	"github.com/valyala/fasthttp"
	"github.com/wI2L/jettison"
)

var (
	conn, err = net.ListenPacket("udp4", ":6969")
	Client    = wrc.NewDebug(conn)
	logger    = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: "FASTHTTP",
	})
)

var upgrader = websocket.FastHTTPUpgrader{
	ReadBufferSize:    4096,
	WriteBufferSize:   4096,
	EnableCompression: true,

	CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
		return true
	},
}

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

func RegisterRoutes() (r router.Router) {
	r = *router.New()
	r.GET("/api/data", jsonHandler)
	r.GET("/api/events", sseHandler)
	r.GET("/api/ws", wsHandler)
	r.GET("/api/avg", avgSpeedHandler)

	printRoutes(&r)
	return r
}

func printRoutes(r *router.Router) {
	routes := r.List()
	for method, paths := range routes {
		for _, path := range paths {
			logger.Info("Registered Route: ", "METHOD", method, "PATH", path)
		}
	}
}

func sseHandler(c *fasthttp.RequestCtx) {
	headers := &c.Response.Header
	headers.Add("Content-Type", "text/event-stream")
	headers.Add("Cache-Control", "no-cache")
	headers.Add("Connection", "keep-alive")
	headers.Add("Transfer-Encoding", "chunked")

	c.SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		data, err := Client.Last()
		if err != nil {
			log.Error(err)
		}
		logger.Info("SSE Client connected!")
		for {

			fmt.Fprintf(w, "event: currentTime\ndata: %s\n\n", data.CurrentStageTime())

			err := w.Flush()
			if err != nil {
				logger.Errorf("Error while flushing: %v. Closing http connection.\n", err)

				break
			}
			fmt.Fprintf(w, "event: inGameTime\ndata: %s\n\n", data.InGameTime())

			err = w.Flush()
			if err != nil {
				logger.Errorf("Error while flushing: %v. Closing http connection.\n", err)

				break
			}
			time.Sleep(100 * time.Millisecond)
		}
	}))
}

func wsHandler(ctx *fasthttp.RequestCtx) {
	err := upgrader.Upgrade(ctx, func(conn *websocket.Conn) {
		defer conn.Close()

		for {

			p, err := Client.Last()
			if err != nil {
				logger.Error(err)
			}
			p.PacketUID++
			err = conn.WriteJSON(p)
			if err != nil {
				if err == websocket.ErrNilConn {
					logger.Error(err)
					conn.Close()
					return
				}
				conn.Close()
				return
			}
			time.Sleep(100 * time.Millisecond)
		}
	})
	if err != nil {
		logger.Error(err)
	}
	logger.Warn("‼️ Websocket client disconnected.")
}

func jsonHandler(c *fasthttp.RequestCtx) {
	p, err := Client.Last()
	if err != nil {
		logger.Error(err)
	}

	data, err := p.ToJSON()
	if err != nil {
		logger.Error(err)
	}
	_, err = c.Write(data)
	if err != nil {
		logger.Error(err)
	}
}

func avgSpeedHandler(c *fasthttp.RequestCtx) {
	s, err := Client.AverageSpeedKmph()
	if err != nil {
		logger.Error(err)
	}
	d, err := jettison.Marshal(map[string]any{
		"speed":   s,
		"samples": Client.Size(),
	})
	if err != nil {
		logger.Error(err)
	}
	_, err = c.Write(d)
	if err != nil {
		logger.Error(err)
	}
}
