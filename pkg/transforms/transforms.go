package transforms

// how do we handle structured vs unstructured.
type Transform interface {
	Encode(event any) []byte
}
