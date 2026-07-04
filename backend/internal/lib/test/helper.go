package testhelper

// Ptr returns pointer by value.
func Ptr[T any](v T) *T {
	return &v
}
