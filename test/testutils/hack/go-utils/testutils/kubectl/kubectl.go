package kubectl

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/rotisserie/eris"
	"github.com/solo-io/go-utils/log"
	"github.com/solo-io/go-utils/threadsafe"
)

type Params struct {
	Stdin          io.Reader
	Stdout, Stderr io.Writer
	Env            []string
	Args           []string
}

func DeleteCrd(ctx context.Context, crd string, params Params) error {
	params.Args = []string{"delete", "crd", crd}
	return Kubectl(ctx, params)
}

func NewParams(args ...string) Params {
	p := Params{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    os.Environ(),
		Args:   args,
	}

	// disable DEBUG=1 from getting through to kube
	for i, pair := range p.Env {
		if strings.HasPrefix(pair, "DEBUG") {
			p.Env = append(p.Env[:i], p.Env[i+1:]...)
			break
		}
	}

	return p
}

func kubectl(ctx context.Context, params Params) *exec.Cmd {
	cmd := exec.CommandContext(ctx, "kubectl", params.Args...)
	cmd.Env = params.Env
	cmd.Stdin = params.Stdin
	cmd.Stdout = params.Stdout
	cmd.Stderr = params.Stderr

	return cmd
}

func Kubectl(ctx context.Context, params Params) error {
	cmd := kubectl(ctx, params)
	log.Debugf("running: %s", strings.Join(cmd.Args, " "))
	return cmd.Run()
}

func KubectlOut(ctx context.Context, params Params) (string, error) {
	// because we are using CombinedOutput we need to set the Stdout and Stderr to nil
	params.Stderr = nil
	params.Stdout = nil
	cmd := kubectl(ctx, params)
	log.Debugf("running: %s", strings.Join(cmd.Args, " "))
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = fmt.Errorf("%s (%v)", out, err)
	}
	return string(out), err
}

func KubectlOutAsync(ctx context.Context, params Params) (io.Reader, error) {
	cmd := kubectl(ctx, params)
	buf := &threadsafe.Buffer{}
	cmd.Stdout = io.MultiWriter(cmd.Stdout, buf)
	cmd.Stderr = io.MultiWriter(cmd.Stderr, buf)

	log.Debugf("async running: %s", strings.Join(cmd.Args, " "))
	err := cmd.Start()
	if err != nil {
		err = fmt.Errorf("%s (%v)", buf.Bytes(), err)
	}

	return buf, err
}

func KubectlOutChan(ctx context.Context, params Params) (<-chan io.Reader, error) {
	cmd := kubectl(ctx, params)
	buf := &threadsafe.Buffer{}
	cmd.Stdout = io.MultiWriter(cmd.Stdout, buf)
	cmd.Stderr = io.MultiWriter(cmd.Stderr, buf)

	log.Debugf("async running: %s", strings.Join(cmd.Args, " "))
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	resultChan := make(chan io.Reader)
	go func() {
		for {
			select {
			case <-time.After(time.Second):
				select {
				case resultChan <- buf:
					continue
				case <-ctx.Done():
					return
				default:
					continue
				}
			}
		}
	}()

	return resultChan, err
}

// WaitPodsRunning waits for all pods to be running
func WaitPodsRunning(ctx context.Context, interval time.Duration, namespace string, params Params, labels ...string) error {
	finished := func(output string) bool {
		return strings.Contains(output, "Running") || strings.Contains(output, "ContainerCreating")
	}
	for _, label := range labels {
		if err := WaitPodStatus(ctx, interval, namespace, label, "Running or ContainerCreating", finished, params); err != nil {
			return err
		}
	}
	finished = func(output string) bool {
		return strings.Contains(output, "Running")
	}
	for _, label := range labels {
		if err := WaitPodStatus(ctx, interval, namespace, label, "Running", finished, params); err != nil {
			return err
		}
	}
	return nil
}

func WaitPodStatus(ctx context.Context, interval time.Duration, namespace, label, status string, finished func(output string) bool, params Params) error {
	tick := time.Tick(interval)
	deadline, _ := ctx.Deadline()
	log.Debugf("waiting till %v for pod %v to be %v...", deadline, label, status)
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for %v to be %v", label, status)
		case <-tick:
			out, err := KubectlOut(ctx, Params{
				Stdin:  params.Stdin,
				Stdout: params.Stdout,
				Stderr: params.Stderr,
				Env:    params.Env,
				Args:   []string{"get", "pod", "-l", label, "-n", namespace},
			})
			if err != nil {
				return fmt.Errorf("failed getting pod: %v", err)
			}
			if strings.Contains(out, "CrashLoopBackOff") {
				out = KubeLogs(ctx, label, params)
				return eris.Errorf("%v in crash loop with logs %v", label, out)
			}
			if strings.Contains(out, "ErrImagePull") || strings.Contains(out, "ImagePullBackOff") {
				out, _ = KubectlOut(ctx, Params{
					Stdin:  params.Stdin,
					Stdout: params.Stdout,
					Stderr: params.Stderr,
					Env:    params.Env,
					Args:   []string{"describe", "pod", "-l", label},
				})
				return eris.Errorf("%v in ErrImagePull with description %v", label, out)
			}
			if finished(out) {
				return nil
			}
		}
	}
}

func KubeLogs(ctx context.Context, label string, params Params) string {
	params.Args = []string{"logs", "-l", label}
	out, err := KubectlOut(ctx, params)
	if err != nil {
		out = err.Error()
	}
	return out
}
