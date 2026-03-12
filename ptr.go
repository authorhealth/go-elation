package elation

// Deprecated: use new() instead of Ptr.
// Will be removed soon after upgrading to Go version 1.26
func Ptr[T any](val T) *T {
	return &val
}
