package redis

import (
	"context"
	"fmt"
	"os/exec"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

type Instance struct {
	startCmd *exec.Cmd
	port     uint32
	address  string

	session *gexec.Session
}

func (i *Instance) Run(ctx context.Context) {
	go func() {
		<-ctx.Done()
		i.Clean()
	}()

	var err error
	Eventually(func(g Gomega) {
		i.session, err = gexec.Start(i.startCmd, GinkgoWriter, GinkgoWriter)
		g.Expect(err).NotTo(HaveOccurred(), "should be able to start redis")
	}, "5s", "1s").Should(Succeed())
	Eventually(i.session.Out, "5s").Should(gbytes.Say("Ready to accept connections"))
}

func (i *Instance) Clean() {
	i.session.Terminate().Wait("1s")
	GinkgoWriter.Println("Redis instance successfully destroyed")
}

func (i *Instance) Port() uint32 {
	return i.port
}

func (i *Instance) Address() string {
	return i.address
}

func (i *Instance) Url() string {
	return fmt.Sprintf("%s:%d", i.Address(), i.Port())
}
