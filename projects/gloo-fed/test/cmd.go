package test

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	. "github.com/onsi/gomega"
	"github.com/solo-io/go-utils/log"
)

// This is a version of Curl that is used for tests running on kind clusters, and do not depend on test runners.
type CurlOpts struct {
	Protocol      string
	Path          string
	Method        string
	Host          string
	Service       string
	CaFile        string
	Body          string
	Headers       map[string]string
	Port          int
	ReturnHeaders bool
}

func CurlEventuallyShouldRespond(opts CurlOpts, substr string, ginkgoOffset int, timeout ...time.Duration) {
	t := time.Second * 20
	if len(timeout) > 0 {
		t = timeout[0]
	}
	EventuallyWithOffset(ginkgoOffset, func() string {
		res, err := Curl(opts)
		if err != nil {
			res = err.Error()
		}
		log.GreyPrintf("running: %v\nwant %v\nhave: %s", opts, substr, res)
		if strings.Contains(res, substr) {
			log.GreyPrintf("success: %v", res)
		}
		return res
	}, t, "5s").Should(ContainSubstring(substr))
}

func Curl(opts CurlOpts) (string, error) {
	args := []string{"-v", "--connect-timeout", "10", "--max-time", "10"}

	if opts.ReturnHeaders {
		args = append(args, "-I")
	}

	if opts.Method != "GET" && opts.Method != "" {
		args = append(args, "-X"+opts.Method)
	}
	if opts.Host != "" {
		args = append(args, "-H", "Host: "+opts.Host)
	}
	if opts.CaFile != "" {
		args = append(args, "--cacert", opts.CaFile)
	}
	if opts.Body != "" {
		args = append(args, "-H", "Content-Type: application/json")
		args = append(args, "-d", opts.Body)
	}
	for h, v := range opts.Headers {
		args = append(args, "-H", fmt.Sprintf("%v: %v", h, v))
	}
	port := opts.Port
	if port == 0 {
		port = 8080
	}
	protocol := opts.Protocol
	if protocol == "" {
		protocol = "http"
	}
	service := opts.Service
	if service == "" {
		service = "test-ingress"
	}
	args = append(args, fmt.Sprintf("%v://%s:%v%s", protocol, service, port, opts.Path))
	log.Debugf("running: curl %v", strings.Join(args, " "))
	cmd := exec.Command("curl", args...)
	cmd.Env = os.Environ()
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%s (%v)", out, err)
	}
	return string(out), err
}
