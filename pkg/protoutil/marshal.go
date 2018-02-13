package protoutil

import (
	"bytes"
	"encoding/json"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/gogo/protobuf/types"
)

var jsonpbMarshaler = &jsonpb.Marshaler{}

func MarshalStruct(m interface{}) (*types.Struct, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	var pb types.Struct
	err = jsonpb.UnmarshalString(string(data), &pb)
	return &pb, err
}

func UnmarshalStruct(structuredData *types.Struct, into interface{}) error {
	strData, err := jsonpbMarshaler.MarshalToString(structuredData)
	if err != nil {
		return err
	}
	data := []byte(strData)
	return json.Unmarshal(data, into)
}

func Unmarshal(data []byte, into proto.Message) error {
	return jsonpb.Unmarshal(bytes.NewBuffer(data), into)
}

func Marshal(pb proto.Message) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := jsonpbMarshaler.Marshal(buf, pb)
	return buf.Bytes(), err
}
