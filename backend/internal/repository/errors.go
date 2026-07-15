package repository

import "errors"

var (
	// ErrNotEnoughStock is returned when stock quantity is insufficient.
	ErrNotEnoughStock = errors.New("not enough stocks")
	// ErrNotFound  is returned when items not found.
	ErrNotFound = errors.New("not found")
)
