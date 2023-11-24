package wrc

import (
	"net"
	"sync"

	"github.com/charmbracelet/log"
)

// Client is used to manage incoming UDP data.
type Client struct {
	ch       chan Packet
	mu       *sync.Mutex
	data     *FIFOQueue[Packet]
	queueLen int
}

// NewWrcClient returns a new client.
func NewWrcClient() *Client {
	return &Client{
		make(chan Packet),
		&sync.Mutex{},
		NewFIFOQueue[Packet](),
		1000,
	}
}

// Pushes new data to the clients queue.
func (w *Client) pushData(p Packet) error {
	w.mu.Lock()
	if w.data.Len() >= w.queueLen {
		return ErrQueueFull
	}
	w.data.Enqueue(p)
	w.mu.Unlock()
	return nil
}

// LatestData returns the latest data from the internal queue.
func (w *Client) LatestData(p Packet) (Packet, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.data.Dequeue()
}

// AverageSpeed returns your average speed based on all 'VehicleSpeed' values in the internal queue.
func (w *Client) AverageSpeed() (float32, error) {
	var speed float32
	data, err := w.data.Data()
	if err != nil {
		return -1, err
	}
	length := w.data.Len()
	for _, s := range data {
		speed += s.VehicleSpeed
	}
	return float32(speed / float32(length)), nil

}

// Run starts the UDP server and begins to listen for incoming packets.
func (w *Client) Run(conn net.PacketConn) error {
	log.WithPrefix("UDP").Info("Started listening for packets.", "address", conn.LocalAddr().String())
	go ListenForPacket(conn, w.ch)

	for p := range w.ch {
		err := w.pushData(p)
		if err != nil {
			return err
		}
	}
	return nil
}

// SetMaxSize of the internal queue. Affects accuracy of AverageSpeed etc.
func (w *Client) SetMaxSize(s int) {
	if w.data.Len() > s {
		w.mu.Lock()
		w.data.Resize(s)
		w.mu.Unlock()
	}
	w.queueLen = s
}
