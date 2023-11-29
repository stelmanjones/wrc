package wrc

import "sync"

type WrcDataStore struct {
	mu       *sync.RWMutex
	store    []*Packet
	max_size int
}

func NewWrcDataStore(d []*Packet) *WrcDataStore {
	p := NewPacket()
	d = append(d, p)
	return &WrcDataStore{
		&sync.RWMutex{},
		d,
		1000,
	}
}

// Pushes new data to the store.
func (w *WrcDataStore) Push(p *Packet) error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if p == nil {
		return ErrNoData
	}
	if len(w.store) >= w.max_size {
		w.store = append(w.store[1:], p)
		return nil
	}
	w.store = append(w.store, p)
	return nil
}

func (w *WrcDataStore) Clear() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.store = make([]*Packet, w.max_size)
	return nil
}

// Last returns the latest packet from the store.
func (w *WrcDataStore) Last() (*Packet, error) {
	w.mu.RLock()
	defer w.mu.RUnlock()
	if len(w.store) == 0 {
		return nil, ErrStoreEmpty
	}

	return w.store[len(w.store)-1], nil
}

// SetSize changes the maixmum number of elements in the store.
func (w *WrcDataStore) SetSize(s int) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if s > len(w.store) {
		w.store = w.store[0:s]
	}
}

// Size returns the current number of elements in the store.
func (w *WrcDataStore) Size() int {
	w.mu.RLock()
	defer w.mu.RUnlock()
	return len(w.store)
}
