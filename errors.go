package wrc

import "errors"

var (
	// ErrNoData is an error that represents a nil data scenario.
	ErrNoData     = errors.New("data is equal to nil")
	
	// ErrStoreEmpty is an error that represents an empty datastore scenario.
	ErrStoreEmpty = errors.New("the datastore is empty")

	// ErrStoreUpdate gets returned if there was a problem updating the
	//internal datastore.
	ErrStoreUpdate = errors.New("there was an error updating the internal store")

	// ErrUDPData represents an udp error.
	ErrUDPData = errors.New("there was an error getting udp data")
)
