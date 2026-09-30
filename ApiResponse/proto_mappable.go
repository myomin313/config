package ApiResponse

// ProtoMappable is implemented by any resource
// that can be sent as protobuf Struct data.
type ProtoMappable interface {
	ToMap() map[string]interface{}
}
