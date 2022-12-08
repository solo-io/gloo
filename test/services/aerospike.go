package services

import (
	ae "github.com/aerospike/aerospike-client-go"
)

const (
	AerospikeDbContainerName = "aerospike"
	AerospikePort            = 3000
)

func RunAerospikeContainer() error {
	image := "aerospike/aerospike-server:6.2.0.0"
	args := []string{"-d", "--rm",
		// this will expose all the ports on aerospike to localhost
		"-p", "3000-3002:3000-3002",
		"--net", GetContainerNetwork(),
		image,
	}
	return RunContainer(AerospikeDbContainerName, args)
}

func ConfigureAerospike() error {
	args := []string{
		// Have to set the nsup-period so that the container will work with the rate limiter
		// if the nsup-period is not set, then a request from the rate limiter will be blocked
		// because TTL is used in the rate limiter
		"asinfo", "-v", "set-config:context=namespace;id=test;nsup-period=1",
	}
	_, err := ExecOnContainer(AerospikeDbContainerName, args)
	return err
}

func AerospikeIsHealthy(address string, port int) bool {
	args := []string{
		"asinfo", "-v", "status",
	}
	out, err := ExecOnContainer(AerospikeDbContainerName, args)
	if err != nil {
		return false
	}
	if len(out) > 0 && string(out) == "ok\n" {
		// if we can connect to the client, then we are good
		_, err = ae.NewClient(address, port)
		return err == nil
	}
	return false
}

func GetAerospikeHost() string {
	return GetDockerHost(AerospikeDbContainerName)
}
