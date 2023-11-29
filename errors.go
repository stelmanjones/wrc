package wrc

import "errors"

var (
	ErrNoData     = errors.New("data is equal to nil")
	ErrStoreEmpty = errors.New("the datastore is empty")
)
