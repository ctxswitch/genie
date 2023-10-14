package transforms

// TODO: placeholders for transforms.  These will take structured (and in cases unstructured) data
// and transform it into other formats/protocols.  I intend this to be used to transform data into
// other formats like protobuf, msgpack, etc.
type Transform interface {
	Encode(event any) []byte
}
