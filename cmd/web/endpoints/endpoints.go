package endpoints

import (
	"embed"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
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




