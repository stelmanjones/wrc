package wrc

import (
	"net"

	"github.com/charmbracelet/log"
)

// Client is used to manage incoming UDP data.
type Client struct {
	*WrcDataStore
	ch chan Packet
}

// NewWrcClient returns a new client.
func NewWrcClient() *Client {
	return &Client{
		NewWrcDataStore(make([]*Packet, 1000)),
		make(chan Packet, 1000),
	}
}

// AverageSpeed returns your average speed based on all 'VehicleSpeed' values in the store.
func (w *Client) AverageSpeed() (float32, error) {
	var speed float32
	w.mu.RLock()
	defer w.mu.RUnlock()
	for _, s := range w.store {
		speed += s.VehicleSpeed
	}
	length := w.Size()
	return float32(speed / float32(length)), nil
}

// Run starts the UDP server and begins to listen for incoming packets.
func (w *Client) Run(conn net.PacketConn) error {
	log.WithPrefix("UDP").Info("Started listening for packets.", "address", conn.LocalAddr().String())
	go ListenForPacket(conn, w.ch)

	for p := range w.ch {
		err := w.Push(&p)
		if err != nil {
			return err
		}
	}
	return nil
}
