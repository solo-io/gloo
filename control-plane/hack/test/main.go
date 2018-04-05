package main

import (
	"encoding/json"

	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/types"
	"github.com/solo-io/gloo/pkg/log"
)

func MapToStruct(m map[string]interface{}) (*types.Struct, error) {
	data, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	var pb types.Struct
	err = jsonpb.UnmarshalString(string(data), &pb)
	return &pb, err
}

func main() {
	foo := map[string]interface{}{
		"hello": "friend",
		"how": struct {
			Are string `json:"are"`
			You string `json:"you"`
		}{
			Are: "are",
			You: "you",
		},
	}
	pb, err := MapToStruct(foo)
	if err != nil {
		panic(err)
	}
	log.Printf("%v", pb)
}
