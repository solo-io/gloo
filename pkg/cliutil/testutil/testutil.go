package testutil

import (
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	expect "github.com/Netflix/go-expect"
	"github.com/hinshun/vt10x"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/solo-io/gloo/pkg/cliutil"
	"gopkg.in/AlecAivazis/survey.v1/terminal"
)

func Stdio(c *expect.Console) terminal.Stdio {
	return terminal.Stdio{c.Tty(), c.Tty(), c.Tty()}
}

// timeout is the max execution time for testCli()
func ExpectInteractive(userInput func(*Console), testCli func(), timeout *time.Duration) {
	c, state, err := vt10x.NewVT10XConsole()
	Expect(err).NotTo(HaveOccurred())
	defer c.Close()
	cliutil.UseStdio(Stdio(c))
	// Dump the terminal's screen.
	defer func() { GinkgoWriter.Write([]byte(expect.StripTrailingEmptyLines(state.String()))) }()

	runId := rand.Int()

	// doneC represents when we're done with the console, indicated by closing the channel
	doneC := make(chan struct{})
	go func() {
		defer func() {
			fmt.Println("ARIANA", time.Now().Format("15:04:05.999999999"), runId, "ExpectInteractive: running GinkgoRecover()")
			GinkgoRecover()
			close(doneC)
			fmt.Println("ARIANA", time.Now().Format("15:04:05.999999999"), runId, "ExpectInteractive: closed doneC")
		}()

		userInput(&Console{console: c})
		fmt.Println("ARIANA", time.Now().Format("15:04:05.999999999"), runId, "ExpectInteractive: userInput func returned; console EOF")
	}()

	//	time.Sleep(time.Hour)
	go func() {
		defer GinkgoRecover()

		testCli()

		fmt.Println("ARIANA", time.Now().Format("15:04:05.999999999"), runId, "ExpectInteractive: finished testCli()")
		// Close the slave end of the pty, and read the remaining bytes from the master end.
		c.Tty().Close()
		fmt.Println("ARIANA", time.Now().Format("15:04:05.999999999"), runId, "ExpectInteractive: closed console")
		<-doneC
		fmt.Println("ARIANA", time.Now().Format("15:04:05.999999999"), runId, "ExpectInteractive: received on doneC")
	}()

	fmt.Println("ARIANA", time.Now().Format("15:04:05.999999999"), runId, "ExpectInteractive: sent off goroutines")

	after := 10 * time.Second
	if timeout != nil {
		after = *timeout
	}
	select {
	case <-time.After(after):
		fmt.Println("ARIANA", time.Now().Format("15:04:05.999999999"), runId, "ExpectInteractive: test timed out")
		c.Tty().Close()
		Fail("test timed out")
	case <-doneC:
		fmt.Println("ARIANA", time.Now().Format("15:04:05.999999999"), runId, "ExpectInteractive: final received on doneC")
	}

	fmt.Println("ARIANA", time.Now().Format("15:04:05.999999999"), runId, "ExpectInteractive: exiting")
}

type Console struct {
	console *expect.Console
}

func (c *Console) ExpectString(s string) string {
	ret, err := c.console.ExpectString(s)
	Expect(err).NotTo(HaveOccurred())
	return ret
}

func (c *Console) PressDown() {
	// These codes are covered here: https://en.wikipedia.org/wiki/ANSI_escape_code
	// see "Escape sequences" and "CSI sequences"
	// 27 = Escape
	// Alternatively, you can use the values written here: gopkg.in/AlecAivazis/survey.v1/terminal/sequences.go
	// But I used the CSI as I seems to be more standard

	_, err := c.console.Write([]byte{27, '[', 'B'})
	Expect(err).NotTo(HaveOccurred())
}

func (c *Console) Esc() {
	// I grabbed this value from here: gopkg.in/AlecAivazis/survey.v1/terminal/sequences.go
	// Originally I tried to use escape codes (https://en.wikipedia.org/wiki/ANSI_escape_code)
	// but it didnt work
	_, err := c.console.Write([]byte{27})
	Expect(err).NotTo(HaveOccurred())
}

func (c *Console) SendLine(s string) int {
	ret, err := c.console.SendLine(s)
	Expect(err).NotTo(HaveOccurred())
	return ret
}

func (c *Console) ExpectEOF() string {
	ret, err := c.console.ExpectEOF()
	Expect(err).NotTo(HaveOccurred())
	fmt.Println("ARIANA", time.Now().Format("15:04:05.999999999"), "CONSOLE EOF", ret)
	return ret
}

type MockKubectl struct {
	Expected        []string
	Next            int
	StdoutLines     []string
	StdoutLineIndex int
}

func NewMockKubectl(cmds []string, stdoutLines []string) *MockKubectl {
	return &MockKubectl{
		Expected:    cmds,
		Next:        0,
		StdoutLines: stdoutLines,
	}
}

func (k *MockKubectl) Kubectl(stdin io.Reader, args ...string) error {
	// If this fails then the CLI tried to run commands we didn't account for in the mock
	Expect(k.Next).To(BeNumerically("<", len(k.Expected)))
	Expect(stdin).To(BeNil())
	cmd := strings.Join(args, " ")
	Expect(cmd).To(BeEquivalentTo(k.Expected[k.Next]))
	k.Next = k.Next + 1
	return nil
}

func (k *MockKubectl) KubectlOut(stdin io.Reader, args ...string) ([]byte, error) {
	Expect(k.Next).To(BeNumerically("<", len(k.Expected)), "MockKubectl did not have a next command for KubectlOut")
	Expect(stdin).To(BeNil(), "Should have passed nil to MockKubectl.KubectlOut")
	cmd := strings.Join(args, " ")
	Expect(cmd).To(BeEquivalentTo(k.Expected[k.Next]), "Wrong next command for MockKubectl.KubectlOut")
	k.Next = k.Next + 1
	Expect(k.StdoutLineIndex).To(BeNumerically("<", len(k.StdoutLines)), "Mock kubectl has run out of stdout lines on command "+cmd)
	stdOutLine := k.StdoutLines[k.StdoutLineIndex]
	k.StdoutLineIndex = k.StdoutLineIndex + 1
	return []byte(stdOutLine), nil
}
