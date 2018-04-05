package consul

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/consul/api"
	"github.com/solo-io/gloo-storage/dependencies"
)

func key(rootPath, itemName string) string {
	return rootPath + "/" + itemName
}

func toKVPair(rootPath string, file *dependencies.File) *api.KVPair {
	var modifyIndex uint64
	if file.ResourceVersion != "" {
		if i, err := strconv.Atoi(file.ResourceVersion); err == nil {
			modifyIndex = uint64(i)
		}
	}
	return &api.KVPair{
		Key:         key(rootPath, file.Ref),
		Value:       file.Contents,
		ModifyIndex: modifyIndex,
	}
}

func setResourceVersion(file *dependencies.File, p *api.KVPair) {
	resourceVersion := fmt.Sprintf("%v", p.ModifyIndex)
	file.ResourceVersion = resourceVersion
}

func fileFromKVPair(rootPath string, p *api.KVPair) *dependencies.File {
	ref := strings.TrimPrefix(p.Key, rootPath+"/")
	return &dependencies.File{
		Ref:             ref,
		Contents:        p.Value,
		ResourceVersion: fmt.Sprintf("%v", p.ModifyIndex),
	}
}
