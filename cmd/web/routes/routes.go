package routes

import (
	"bufio"
	"embed"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/valyala/fasthttp"
)


func CorsHandler() func(*fiber.Ctx) error {
	return cors.New(cors.Config{
		AllowOriginsFunc: func(origin string) bool {
			return os.Getenv("ENV") == "debug"
		},
	})
}

func EmbedFSHandler(fs *embed.FS) func(*fiber.Ctx) error {
	return filesystem.New(filesystem.Config{
		Root:       http.FS(fs),
		PathPrefix: "assets",
		Browse:     true,
	})
}

func RootHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		return c.Render("main", fiber.Map{}, "layouts/shell")
	}
}

func SSEHandler() func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		c.Set("Content-Type", "text/event-stream")
		c.Set("Cache-Control", "no-cache")
		c.Set("Connection", "keep-alive")
		c.Set("Transfer-Encoding", "chunked")

		c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
			var i int
			log.Info("SSE Client connected!")
			for {
				i++
				msg := fmt.Sprintf("%d", i)
				fmt.Fprintf(w, "data: %s\n\n", msg)
				log.Debug("SSE", "count", i)

				err := w.Flush()
				if err != nil {
					// Refreshing page in web browser will establish a new
					// SSE connection, but only (the last) one is alive, so
					// dead connections must be closed here.
					log.Errorf("Error while flushing: %v. Closing http connection.\n", err)

					break
				}
				time.Sleep(1 * time.Second)
			}
		}))
		return nil
	}

}