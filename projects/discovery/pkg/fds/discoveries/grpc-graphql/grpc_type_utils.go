package grpc

import (
	"github.com/jhump/protoreflect/desc"
)

const (
	GoogleProtobufPackage = "google.protobuf"
)

// Certain types in the `google.protobuf` package have custom encodings in JSON.
// as seen in the JSON encoder in envoy:
// https://github.com/protocolbuffers/protobuf/blob/e7cc1aa970a053d366d4f7faa2a22ecb356292c4/src/google/protobuf/util/internal/protostream_objectsource.cc#L675-L710
func TranslateGoogleProtobufWrapperTypes(descriptor desc.Descriptor) string {
	if descriptor.GetFile().GetPackage() == GoogleProtobufPackage {
		var newFieldType string
		// todo - support 'Struct', 'ListValue', and 'Any' which translate to JSON objects/arrays with any field names
		// which GraphQL does not currently support -- this requires a dataplane change to support a JSON type
		switch descriptor.GetName() {
		case "DoubleValue", "FloatValue":
			newFieldType = GRAPHQL_FLOAT
		case "Int32Value", "UInt32Value":
			newFieldType = GRAPHQL_INT
		case "BoolValue":
			newFieldType = GRAPHQL_BOOLEAN
		case "Int64Value", "UInt64Value", "BytesValue", "Timestamp", "Duration", "StringValue", "FieldMask":
			newFieldType = GRAPHQL_STRING
		}
		return newFieldType
	}
	return ""
}
