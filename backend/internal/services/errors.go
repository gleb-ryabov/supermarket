package services

import "errors"

var (
	// ErrNotEnoughStock is returned when stock quantity is insufficient.
	ErrNotEnoughStock = errors.New("not enough stocks")
)
