package envoy

import (
	adminv3 "github.com/envoyproxy/go-control-plane/envoy/admin/v3"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"os"
)

// Function to parse Envoy configuration
func parseEnvoyConfig(inputFile string) (*EnvoySnapshot, error) {

	// Read the Envoy configuration file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	// Unmarshal the JSON into the ConfigDump struct
	var configDump adminv3.ConfigDump
	if err := protojson.Unmarshal(data, &configDump); err != nil {
		log.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	envoysnapshot := &EnvoySnapshot{}

	for _, config := range configDump.Configs {
		if config.GetTypeUrl() == "type.googleapis.com/envoy.admin.v3.ListenersConfigDump" {
			if err := config.UnmarshalTo(&envoysnapshot.Listeners); err != nil {
				log.Fatalf("Failed to unmarshal message: %v", err)
			}
		}
		if config.GetTypeUrl() == "type.googleapis.com/envoy.admin.v3.RoutesConfigDump" {
			if err := config.UnmarshalTo(&envoysnapshot.Routes); err != nil {
				log.Fatalf("Failed to unmarshal message: %v", err)
			}
		}
		if config.GetTypeUrl() == "type.googleapis.com/envoy.admin.v3.ClustersConfigDump" {
			if err := config.UnmarshalTo(&envoysnapshot.Clusters); err != nil {
				log.Fatalf("Failed to unmarshal message: %v", err)
			}
		}
	}

	return envoysnapshot, nil
}
