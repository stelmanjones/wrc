package routes

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/fasthttp/router"
	"github.com/stelmanjones/wrc"
	"github.com/valyala/fasthttp"
)

var (
	Client = wrc.NewWrcClient()
	logger = log.NewWithOptions(os.Stderr, log.Options{
		Prefix: "FASTHTTP",
	})
)

func RegisterRoutes() (r router.Router) {
	r = *router.New()
	r.GET("/api/data", jsonHandler)
	r.GET("/api/events", sseHandler)

	printRoutes(&r)
	return r
}

func printRoutes(r *router.Router) {
	routes := r.List()
	for method, paths := range routes {
		for _, path := range paths {
			logger.Info("Registered Route: ", "METHOD", fmt.Sprintf("%s", method), "PATH", path)
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

func jsonHandler(c *fasthttp.RequestCtx) {
	p, err := Client.Last()
	p.PacketUID += 1
	t := time.Now().Format("03:04:05.000")
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
