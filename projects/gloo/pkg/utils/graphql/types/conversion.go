package types

import (
	"github.com/golang/protobuf/proto"
	"github.com/rotisserie/eris"
)

func ConvertGoProtoTypes(inputMessage proto.Message, outputProtoMessage proto.Message) error {
	protoIntermediateBytes, err := proto.Marshal(inputMessage)
	if err != nil {
		return eris.Wrapf(err, "proto message %s cannot be marshalled", inputMessage.String())
	}
	err = proto.Unmarshal(protoIntermediateBytes, outputProtoMessage)
	if err != nil {
		return eris.Wrapf(err, "proto message %s cannot be unmarshalled into proto message %s", inputMessage.String(), outputProtoMessage.String())
	}
	return nil
}
