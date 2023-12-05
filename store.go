package wrc

import "sync"

// DataStore represents a data store for Wrc packets with thread-safe operations.
type DataStore struct {
	mu       *sync.RWMutex
	store    []*Packet
	maxSize int
}
// NewWrcDataStore initializes a new WrcDataStore with the provided packets.
func NewWrcDataStore(d []*Packet) *DataStore {
	p := NewPacket()
	d = append(d, p)
	return &DataStore{
		&sync.RWMutex{},
		d,
		600,
	}
}

// Push new data to the store.
func (w *DataStore) Push(p *Packet) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if p == nil {
		return ErrNoData
	}
	if len(w.store) == w.maxSize {
		w.store = w.store[1:]
		w.store = append(w.store[0:], p)
		return nil
	}
	w.store = append(w.store[0:], p)
	return nil
}
// Clear the internal store.
func (w *DataStore) Clear() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.store = make([]*Packet, 0, w.maxSize)
	return nil
}

// Last returns the latest packet from the store.
func (w *DataStore) Last() (*Packet, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if len(w.store) == 0 {
		return nil, ErrStoreEmpty
	}

	return w.store[len(w.store)-1], nil
}

// SetSize changes the maixmum number of elements in the store.
func (w *DataStore) SetSize(s int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if s > len(w.store) {
		w.store = w.store[0:s]
	}
}

// Size returns the current number of elements in the store.
func (w *DataStore) Size() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.store)
}
