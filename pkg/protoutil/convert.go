package protoutil

import (
	"encoding/json"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/ptypes/struct"
)

func MapToStruct(m map[string]interface{}) (*structpb.Struct, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	var pb structpb.Struct
	err = jsonpb.UnmarshalString(string(data), &pb)
	return &pb, err
}
