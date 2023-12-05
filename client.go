package wrc

import (
	"net"
	"os"

	"github.com/charmbracelet/log"
)

var clogger = log.New(os.Stderr).WithPrefix("CLIENT")

// Client is used to manage incoming UDP data.
type Client struct {
	conn net.PacketConn
	*DataStore
	ch    chan Packet
	Debug bool
}

// New returns a new client.
func New(conn net.PacketConn) *Client {
	clogger.Info("WRC Client initialized! 🏁")

	return &Client{
		conn,
		NewWrcDataStore(make([]*Packet, 0, 600)),
		make(chan Packet, 600),
		false,
	}
}

// NewDebug returns a new client with some extra debug info.
func NewDebug(conn net.PacketConn) *Client {
	clogger.Info("WRC Client initialized! 🏁")

	return &Client{
		conn,
		NewWrcDataStore(make([]*Packet, 0, 600)),
		make(chan Packet, 600),
		true,
	}
}

// AverageSpeedKmph returns your average speed based on all 'VehicleSpeed' values in the store.
func (c *Client) AverageSpeedKmph() (float32, error) {
	var speed float32
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, s := range c.store {
		speed += s.VehicleSpeed
	}
	length := c.Size()
	return float32(speed/float32(length)) * MpsToKmph, nil
}
// AverageSpeedMph returns your average speed based on all 'VehicleSpeed' values in the store.
func (c *Client) AverageSpeedMph() (float32, error) {
	var speed float32
	c.mu.RLock()
	defer c.mu.RUnlock()
	for _, s := range c.store {
		speed += s.VehicleSpeed
	}
	length := len(c.store)
	return float32(speed/float32(length)) * MpsToMph, nil
}

// Run starts the UDP client and begins to listen for incoming packets.
func (c *Client) Run() error {
	clogger.Info("Started listening for packets.", "address", c.conn.LocalAddr().String())
	go Listen(c.conn, c.ch)

	for p := range c.ch {
		err := c.Push(&p)
		if err != nil {
			return err

		}
	}
	clogger.Info("Bye! 👋")
	return nil
}
