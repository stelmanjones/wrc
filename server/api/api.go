package api

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/gofiber/fiber/v2"
	"github.com/stelmanjones/wrc"
)

var (
	logger log.Logger = *log.New(os.Stdout).WithPrefix("API")


	udpAddress string
	apiAddress string
)

func RunHttpServer(f *fiber.App, address string, packet *wrc.ThreadSafePacket) error {


	logger.Fatal(f.Listen(address))

	return nil
}
